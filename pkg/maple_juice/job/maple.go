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
)

type MapleTaskStatus int

const (
	Init MapleTaskStatus = iota + 1
	Running
	Errored
	Failed
)

type MapleTask struct {
	TaskID          string
	SplitInfo       *split.Info
	Request         *idl.RunMapleTaskRequest
	Status          MapleTaskStatus
	AttemptCount    int
	LastAttemptNode string
}

type MapleJobTracker struct {
	req              *idl.ExecuteMapleJobRequest
	rpcClientManager *rpc.ClientManager
	fsClient         *SDFSSDK.SDFSClient
	jobManager       *Manager

	splitsInfo        []*split.Info
	tasks             map[string]*MapleTask
	ongoingTasks      chan *MapleTask
	successResponse   chan *idl.RunMapleTaskResponse
	errChan           chan error
	exeLocateHosts    []string
	tempIntermediates map[string][]string

	resp *idl.ExecuteMapleJobResponse
}

func NewMapleJobTracker(req *idl.ExecuteMapleJobRequest, rpcClientManager *rpc.ClientManager,
	fsClient *SDFSSDK.SDFSClient, jobManager *Manager) *MapleJobTracker {
	tracker := MapleJobTracker{
		req:              req,
		rpcClientManager: rpcClientManager,
		fsClient:         fsClient,
		jobManager:       jobManager,

		splitsInfo:        nil,
		tasks:             make(map[string]*MapleTask),
		ongoingTasks:      nil,
		successResponse:   nil,
		errChan:           make(chan error),
		exeLocateHosts:    nil,
		tempIntermediates: make(map[string][]string),
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
	t.ongoingTasks = make(chan *MapleTask, len(t.splitsInfo))
	t.successResponse = make(chan *idl.RunMapleTaskResponse, len(t.splitsInfo))
	for idx, splitInfo := range t.splitsInfo {
		taskID := fmt.Sprintf("task%d", idx+1)

		req := idl.RunMapleTaskRequest{
			ExeName:                    t.req.GetExeName(),
			IntermediateFilenamePrefix: t.req.GetIntermediateFilenamePrefix(),
			InputSplits:                []*idl.InputSplit{splitInfo.Split},
			AttemptId:                  "",
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
		t.ongoingTasks <- &task
	}

	fmt.Println("generated tasks:----------")
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
		case task := <-t.ongoingTasks:
			{
				go t.runTasks(task)
			}
		case resp := <-t.successResponse:
			{
				// response of failed tasks will be processed in runTask func and reschedule it
				t.processResp(resp)
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
	for intermediateName, tmps := range t.tempIntermediates {
		err := t.fsClient.TouchFile(intermediateName)
		if err != nil {
			return fmt.Errorf("can not create intermediate:%w", err)
		}
		logutil.Logger.Debugf("create intermediate (%s)", intermediateName)

		err = t.fsClient.MergeFiles(intermediateName, tmps, true, false)
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
		Code:              idl.StatusCode_Success,
		IntermediateFiles: intermediates,
	}

	t.resp = &resp

	return nil
}

func (t *MapleJobTracker) runTasks(task *MapleTask) {
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
		t.ongoingTasks <- task
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
			if available == nil || len(available) == 0 {
				return "", fmt.Errorf("can not get any available node")
			}
			selectedHost = available[rand.Intn(len(available))]
		}
	}

	return selectedHost, nil
}

func (t *MapleJobTracker) processResp(resp *idl.RunMapleTaskResponse) {
	for _, tmpIntermediate := range resp.GetTmpIntermediateFiles() {
		parts := strings.Split(tmpIntermediate, "-")
		suffix := parts[len(parts)-1]
		intermediate := strings.TrimSuffix(tmpIntermediate, "-"+suffix)

		t.tempIntermediates[intermediate] = append(t.tempIntermediates[intermediate], tmpIntermediate)
	}
}
