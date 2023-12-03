package node

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/gob"
	"fmt"
	"google.golang.org/protobuf/proto"
	"os"
	"os/exec"
)

func (h *ManageServiceHandler) RunMapleTask(ctx context.Context, request *idl.RunMapleTaskRequest) (*idl.RunMapleTaskResponse, error) {
	return NewRunMapleTaskHandler(ctx, request, h.fsClient).Handle()
}

type RunMapleTaskHandler struct {
	ctx      context.Context
	req      *idl.RunMapleTaskRequest
	fsClient *SDFSSDK.SDFSClient

	resp         *idl.RunMapleTaskResponse
	exeStorePath string
}

func NewRunMapleTaskHandler(ctx context.Context, request *idl.RunMapleTaskRequest, client *SDFSSDK.SDFSClient) *RunMapleTaskHandler {
	handler := RunMapleTaskHandler{
		ctx:      ctx,
		req:      request,
		fsClient: client,
	}

	return &handler
}

func (h *RunMapleTaskHandler) Handle() (*idl.RunMapleTaskResponse, error) {
	logutil.Logger.Debugf("start run maple task attempt:%s", h.req.GetAttemptId())
	for _, handleFUnc := range []func() error{
		h.loadExeFile, h.runExeFile, h.processResp,
	} {
		if err := handleFUnc(); err != nil {
			logutil.Logger.Error(err)
			return &idl.RunMapleTaskResponse{
				Code:                 idl.StatusCode_InternalErr,
				TmpIntermediateFiles: nil,
				ErrInfo:              proto.String(err.Error()),
			}, nil
		}
	}

	logutil.Logger.Debugf("execute maple task (%s) sucess, generate intermediates:%v", h.req.GetAttemptId(), h.resp.GetTmpIntermediateFiles())

	return h.resp, nil
}

func (h *RunMapleTaskHandler) loadExeFile() error {
	exeName := h.req.GetExeName()
	h.exeStorePath = DefaultStoreDir + "/" + h.req.GetAttemptId()
	os.MkdirAll(h.exeStorePath, 0777)
	err := h.fsClient.GetFileToLocal(exeName, exeName, h.exeStorePath)
	if err != nil {
		return fmt.Errorf("load exe file failed:%w", err)
	}
	logutil.Logger.Debugf("load exe file (%s) sueccess", exeName)
	return nil
}

func (h *RunMapleTaskHandler) runExeFile() error {
	logutil.Logger.Debugf("start to run exe file (%s)", h.req.GetExeName())
	exeCmd := exec.Command("./"+h.req.GetExeName(), h.req.GetExeArgs()...)
	exeCmd.Dir = h.exeStorePath

	exeStdin, _ := exeCmd.StdinPipe()
	exeStdout, _ := exeCmd.StdoutPipe()

	err := exeCmd.Start()
	if err != nil {
		return fmt.Errorf("can not run exe file:%w", err)
	}

	// forward request to exe process
	encoder := gob.NewEncoder(exeStdin)
	err = encoder.Encode(h.req)
	if err != nil {
		return fmt.Errorf("forward rqeuest to exe failed:%w", err)
	}

	// read response of exe process
	decoder := gob.NewDecoder(exeStdout)
	resp := idl.RunMapleTaskResponse{}
	err = decoder.Decode(&resp)
	if err != nil {
		return fmt.Errorf("get response from exe failed:%w", err)
	}

	h.resp = &resp

	err = exeStdin.Close()
	if err != nil {
		return fmt.Errorf("close exe failed:%w", err)
	}
	err = exeCmd.Wait()
	if err != nil {
		return fmt.Errorf("wait exe cmd finish failed:%w", err)
	}

	return nil
}

func (h *RunMapleTaskHandler) processResp() error {
	if h.resp == nil {
		return fmt.Errorf("get empty response")
	}

	if h.resp.GetCode() != idl.StatusCode_Success {
		return fmt.Errorf("execute failed, err:%s", h.resp.GetErrInfo())
	}

	return nil
}
