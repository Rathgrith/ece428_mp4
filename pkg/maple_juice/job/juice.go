package job

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice/juice/shuffle"
	"ece428_mp4/pkg/rpc"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

type JuiceTask struct {
	TaskID            string
	IntermediateNames []string
	Request           *idl.RunJuiceTaskRequest
	Status            TaskStatus
	AttemptCount      int
	LastAttemptNode   string
	Err               error
}

type JuiceJobTracker struct {
	req              *idl.ExecuteJuiceJobRequest
	rpcClientManager *rpc.ClientManager
	fsClient         *SDFSSDK.SDFSClient
	jobManager       *Manager
	ctx              context.Context

	shuffledIntermediates [][]string
	tasks                 map[string]*JuiceTask
	taskQueue             chan *JuiceTask
	retryQueue            chan *JuiceTask
	successResponse       chan *idl.RunJuiceTaskResponse
	errChan               chan error

	tmpOutputNames []string
}

func NewJuiceJobTracker(ctx context.Context, req *idl.ExecuteJuiceJobRequest, rpcClientManager *rpc.ClientManager,
	fsClient *SDFSSDK.SDFSClient, jobManager *Manager) *JuiceJobTracker {
	tracker := JuiceJobTracker{
		req:              req,
		rpcClientManager: rpcClientManager,
		fsClient:         fsClient,
		jobManager:       jobManager,
		ctx:              ctx,

		shuffledIntermediates: nil,
		tasks:                 make(map[string]*JuiceTask),
		taskQueue:             nil,
		retryQueue:            nil,
		successResponse:       nil,
		errChan:               make(chan error, 1),
		tmpOutputNames:        make([]string, 0),
	}

	return &tracker
}

func (t *JuiceJobTracker) ExecuteJob() (*idl.ExecuteJuiceJobResponse, error) {
	for _, handleFUnc := range []func() error{
		t.shuffleIntermediates, t.generateTasks, t.dispatchAndMonitor, t.mergeTmpOutput,
	} {
		if err := handleFUnc(); err != nil {
			logutil.Logger.Errorf("job failed:%v", err)
			return nil, err
		}
	}

	resp := idl.ExecuteJuiceJobResponse{Code: idl.StatusCode_Success}

	return &resp, nil
}

func (t *JuiceJobTracker) shuffleIntermediates() error {
	shufflerCreateFunc := shuffle.GetShufflerCreateFunc(shuffle.HashShufflerType)
	shuffler := shufflerCreateFunc(t.req.GetIntermediateFilenames(), int(t.req.GetNumMaples()))
	res, err := shuffler.Shuffle()
	if err != nil {
		return fmt.Errorf("can not shuffle input intermediates:%w", err)
	}

	t.shuffledIntermediates = res

	return nil
}

func (t *JuiceJobTracker) generateTasks() error {
	t.taskQueue = make(chan *JuiceTask, len(t.shuffledIntermediates))
	t.successResponse = make(chan *idl.RunJuiceTaskResponse, len(t.shuffledIntermediates))
	t.retryQueue = make(chan *JuiceTask, len(t.shuffledIntermediates))

	for idx, partIntermediates := range t.shuffledIntermediates {
		taskID := fmt.Sprintf("juice_task%d", idx+1)

		req := idl.RunJuiceTaskRequest{
			ExeName:                t.req.GetExeName(),
			InputIntermediateFiles: partIntermediates,
			OutPutFilename:         t.req.GetOutPutFilename(),
			AttemptId:              "",
			ExeArgs:                t.req.GetExeArgs(),
		}

		task := JuiceTask{
			TaskID:            taskID,
			IntermediateNames: partIntermediates,
			Request:           &req,
			Status:            Init,
			AttemptCount:      0,
			LastAttemptNode:   "",
		}

		t.tasks[taskID] = &task
		t.taskQueue <- &task
	}

	fmt.Println("generated juice tasks:----------")
	for _, task := range t.tasks {
		marshaledTask, _ := json.Marshal(task)
		fmt.Println(string(marshaledTask))
	}

	return nil
}

func (t *JuiceJobTracker) dispatchAndMonitor() error {
	remainTasksCount := len(t.tasks)
	for {
		select {
		case task := <-t.taskQueue:
			{
				go t.runTask(t.ctx, task)
			}
		case task := <-t.retryQueue:
			{
				if task.AttemptCount < MaxRetryTime {
					go t.runTask(t.ctx, task)
				} else {
					task.Status = Failed
					err := fmt.Errorf("exceed max retry time, still failed:%w", task.Err)
					logutil.Logger.Error(err)
					return err
				}
			}
		case resp := <-t.successResponse:
			{
				t.processResp(t.ctx, resp)
				remainTasksCount -= 1
				if remainTasksCount == 0 {
					return nil
				}
			}
		case err := <-t.errChan:
			{
				return fmt.Errorf("unable to handle err:%w", err)
			}
		}
	}
}

func (t *JuiceJobTracker) mergeTmpOutput() error {
	outName := t.req.GetOutPutFilename()
	err := t.fsClient.TouchFile(outName)
	if err != nil {
		return fmt.Errorf("can not create output file:%w", err)
	}

	err = t.fsClient.MergeFiles(outName, t.tmpOutputNames, true, false)
	if err != nil {
		return fmt.Errorf("can not merge tmp output files:%w", err)
	}

	logutil.Logger.Debugf("create output file (%s)", t.req.GetOutPutFilename())

	return nil
}

func (t *JuiceJobTracker) runTask(ctx context.Context, task *JuiceTask) {
	taskReq := task.Request

	if len(task.IntermediateNames) == 0 {
		t.successResponse <- &idl.RunJuiceTaskResponse{
			Code:              idl.StatusCode_Success,
			TmpOutputFilename: "",
			ErrInfo:           nil,
		}
		return
	}

	task.AttemptCount += 1
	attemptID := fmt.Sprintf("%s_%d", task.TaskID, task.AttemptCount)
	taskReq.AttemptId = attemptID

	dispatchHost, err := t.selectTargetHost(task)
	if err != nil {
		// job tracker internal error, should not retry
		t.errChan <- fmt.Errorf("can not select host for task (%s):%w", task.TaskID, err)
		return
	}
	task.LastAttemptNode = dispatchHost
	task.Status = Running

	logutil.Logger.Debugf("dispatch task (%s) to node (%s)", task.TaskID, task.LastAttemptNode)

	nodeManagerClient, err := t.rpcClientManager.GetNodeManagerRPCClient(dispatchHost)
	if err != nil {
		t.errChan <- fmt.Errorf("can not get rpc client of host (%s):%w", dispatchHost, err)
		return
	}

	resp, err := nodeManagerClient.RunJuiceTask(context.Background(), taskReq)
	if err != nil {
		// err can only be rpc call err,
		// since when node manager receive request, the run error will be contained in resp
		// just retry and select another node
		task.Err = err
		task.Status = Errored
		t.retryQueue <- task
		return
	}

	if resp.GetCode() != idl.StatusCode_Success {
		t.errChan <- fmt.Errorf("run task err at (%s):%w", dispatchHost, errors.New(resp.GetErrInfo()))
		return
	}

	t.successResponse <- resp
}

func (t *JuiceJobTracker) selectTargetHost(task *JuiceTask) (string, error) {
	available := t.jobManager.GetAvailableHost()
	if len(available) == 0 {
		return "", fmt.Errorf("can not get any available node")
	}
	selectedHost := available[rand.Intn(len(available))]
	return selectedHost, nil
}

func (t *JuiceJobTracker) processResp(ctx context.Context, resp *idl.RunJuiceTaskResponse) {
	if resp.GetTmpOutputFilename() == "" {
		return
	}
	t.tmpOutputNames = append(t.tmpOutputNames, resp.GetTmpOutputFilename())
}
