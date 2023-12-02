package job

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice/maple/split"
	"ece428_mp4/pkg/rpc"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

type MapleTask struct {
	TaskID          string
	SplitInfo       *split.Info
	Request         *idl.RunMapleTaskRequest
	Status          TaskStatus
	AttemptCount    int
	LastAttemptNode string
	Err             error
}

type MapleJobTracker struct {
	req              *idl.ExecuteMapleJobRequest
	rpcClientManager *rpc.ClientManager
	fsClient         *SDFSSDK.SDFSClient
	jobManager       *Manager
	ctx              context.Context

	splitsInfo          []*split.Info
	tasks               map[string]*MapleTask
	taskQueue           chan *MapleTask
	retryQueue          chan *MapleTask
	successResponse     chan *idl.RunMapleTaskResponse
	errChan             chan error
	exeLocateHosts      []string
	tempIntermediates   map[string][]string
	tempIntermediatesMu sync.RWMutex

	resp *idl.ExecuteMapleJobResponse
}

func NewMapleJobTracker(ctx context.Context, req *idl.ExecuteMapleJobRequest, rpcClientManager *rpc.ClientManager,
	fsClient *SDFSSDK.SDFSClient, jobManager *Manager) *MapleJobTracker {
	tracker := MapleJobTracker{
		req:              req,
		rpcClientManager: rpcClientManager,
		fsClient:         fsClient,
		jobManager:       jobManager,
		ctx:              ctx,

		splitsInfo:          nil,
		tasks:               make(map[string]*MapleTask),
		taskQueue:           nil,
		retryQueue:          nil,
		successResponse:     nil,
		errChan:             make(chan error),
		exeLocateHosts:      nil,
		tempIntermediates:   make(map[string][]string),
		tempIntermediatesMu: sync.RWMutex{},
	}

	return &tracker
}

func (t *MapleJobTracker) ExecuteJob() (*idl.ExecuteMapleJobResponse, error) {
	for _, handleFUnc := range []func() error{
		t.splitInputFiles, t.generateTasks, t.dispatchAndMonitor, t.mergeTmpIntermediates, t.generateJobResponse,
	} {
		if err := handleFUnc(); err != nil {
			logutil.Logger.Errorf("job failed:%v", err)
			return nil, err
		}
	}
	return t.resp, nil
}

func (t *MapleJobTracker) splitInputFiles() error {
	// TODO: config splitter
	s := split.NewRawFileSplitter(t.req.GetInputFiles(), t.fsClient, int(t.req.GetNumMaples()))
	splitsInfo, err := s.Split()
	if err != nil {
		return fmt.Errorf("can not spilt input files:%w", err)
	}

	t.splitsInfo = splitsInfo
	return nil
}

func (t *MapleJobTracker) generateTasks() error {
	t.taskQueue = make(chan *MapleTask, len(t.splitsInfo))
	t.successResponse = make(chan *idl.RunMapleTaskResponse, len(t.splitsInfo))
	t.retryQueue = make(chan *MapleTask, len(t.splitsInfo))
	for idx, splitInfo := range t.splitsInfo {
		taskID := fmt.Sprintf("maple_task%d", idx+1)

		req := idl.RunMapleTaskRequest{
			ExeName:                    t.req.GetExeName(),
			IntermediateFilenamePrefix: t.req.GetIntermediateFilenamePrefix(),
			InputSplits:                []*idl.InputSplit{splitInfo.Split},
			AttemptId:                  "",
			ExeArgs:                    t.req.ExeArgs,
		}

		task := MapleTask{
			TaskID:          taskID,
			SplitInfo:       splitInfo,
			Request:         &req,
			Status:          Init,
			AttemptCount:    0,
			LastAttemptNode: "",
		}

		t.tasks[taskID] = &task
		t.taskQueue <- &task
	}

	fmt.Println("generated maple tasks:----------")
	for _, task := range t.tasks {
		marshaledTask, _ := json.Marshal(task)
		fmt.Println(string(marshaledTask))
	}

	return nil
}

func (t *MapleJobTracker) dispatchAndMonitor() error {
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

func (t *MapleJobTracker) mergeTmpIntermediates() error {
	for intermediateName, temp := range t.tempIntermediates {
		err := t.fsClient.TouchFile(intermediateName)
		if err != nil {
			return fmt.Errorf("can not create intermediate:%w", err)
		}
		logutil.Logger.Debugf("create intermediate (%s)", intermediateName)

		err = t.fsClient.MergeFiles(intermediateName, temp, true, false)
		if err != nil {
			return fmt.Errorf("can not merge tmp intermediates:%w", err)
		}
	}
	return nil
}

func (t *MapleJobTracker) generateJobResponse() error {
	intermediates := make([]string, 0, len(t.tempIntermediates))
	for filename := range t.tempIntermediates {
		intermediates = append(intermediates, filename)
	}

	resp := idl.ExecuteMapleJobResponse{
		Code:                  idl.StatusCode_Success,
		IntermediateFilenames: intermediates,
	}

	t.resp = &resp

	return nil
}

func (t *MapleJobTracker) runTask(ctx context.Context, task *MapleTask) {
	req := task.Request

	task.AttemptCount += 1
	attemptID := fmt.Sprintf("%s_%d", task.TaskID, task.AttemptCount)
	req.AttemptId = attemptID

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

	resp, err := nodeManagerClient.RunMapleTask(context.Background(), req)
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

func (t *MapleJobTracker) selectTargetHost(task *MapleTask) (string, error) {
	var preferredHosts []string
	if t.exeLocateHosts == nil {
		hosts, err := t.fsClient.GetFileLocationHosts(task.Request.GetExeName())
		if err != nil {

		} else {
			t.exeLocateHosts = hosts
		}
	}

	preferredHosts = t.exeLocateHosts
	if preferredHosts == nil {
		preferredHosts = task.SplitInfo.Locations
	}

	var selectedHost string
	count := 0
	for {
		selectedHost = preferredHosts[rand.Intn(len(preferredHosts))]
		if t.jobManager.NodeAlive(selectedHost) {
			break
		}
		count += 1
		if count > 2*len(preferredHosts) {
			available := t.jobManager.GetAvailableHost()
			if len(available) == 0 {
				return "", fmt.Errorf("can not get any available node")
			}
			selectedHost = available[rand.Intn(len(available))]
		}
	}

	return selectedHost, nil
}

func (t *MapleJobTracker) processResp(ctx context.Context, resp *idl.RunMapleTaskResponse) {
	t.tempIntermediatesMu.Lock()
	defer t.tempIntermediatesMu.Unlock()
	for _, tmpIntermediate := range resp.GetTmpIntermediateFiles() {
		parts := strings.Split(tmpIntermediate, "-")
		suffix := parts[len(parts)-1]
		intermediate := strings.TrimSuffix(tmpIntermediate, "-"+suffix)

		t.tempIntermediates[intermediate] = append(t.tempIntermediates[intermediate], tmpIntermediate)
	}
}
