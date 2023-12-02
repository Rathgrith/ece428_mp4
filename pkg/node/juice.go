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

func (h *ManageServiceHandler) RunJuiceTask(ctx context.Context, request *idl.RunJuiceTaskRequest) (*idl.RunJuiceTaskResponse, error) {
	return NewRunJuiceTaskHandler(ctx, request, h.fsClient).Handle()
}

type RunJuiceTaskHandler struct {
	ctx      context.Context
	req      *idl.RunJuiceTaskRequest
	fsClient *SDFSSDK.SDFSClient

	resp         *idl.RunJuiceTaskResponse
	exeStorePath string
}

func NewRunJuiceTaskHandler(ctx context.Context, request *idl.RunJuiceTaskRequest, client *SDFSSDK.SDFSClient) *RunJuiceTaskHandler {
	handler := RunJuiceTaskHandler{
		ctx:      ctx,
		req:      request,
		fsClient: client,
	}

	return &handler
}

func (h *RunJuiceTaskHandler) Handle() (*idl.RunJuiceTaskResponse, error) {
	for _, handleFUnc := range []func() error{
		h.loadExeFile, h.runExeFile, h.processResp,
	} {
		if err := handleFUnc(); err != nil {
			logutil.Logger.Error(err)
			return &idl.RunJuiceTaskResponse{
				Code:              idl.StatusCode_InternalErr,
				TmpOutputFilename: "",
				ErrInfo:           proto.String(err.Error()),
			}, nil
		}
	}

	logutil.Logger.Debugf("execute juice task (%s) sucess, generate output:%v", h.req.GetAttemptId(), h.resp.GetTmpOutputFilename())

	return h.resp, nil
}

func (h *RunJuiceTaskHandler) loadExeFile() error {
	exeName := h.req.GetExeName()
	h.exeStorePath = DefaultStoreDir + "/" + h.req.GetAttemptId()
	os.MkdirAll(h.exeStorePath, 0777)
	err := h.fsClient.GetFileToLocal(exeName, exeName, h.exeStorePath)
	if err != nil {
		return fmt.Errorf("load exe file failed:%w", err)
	}
	return nil
}

func (h *RunJuiceTaskHandler) runExeFile() error {
	exeCmd := exec.Command("./" + h.req.GetExeName())
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
	resp := idl.RunJuiceTaskResponse{}
	err = decoder.Decode(&resp)
	if err != nil {
		return fmt.Errorf("get response from exe failed:%w", err)
	}

	h.resp = &resp

	return nil
}

func (h *RunJuiceTaskHandler) processResp() error {
	if h.resp == nil {
		return fmt.Errorf("get empty response")
	}

	if h.resp.GetCode() != idl.StatusCode_Success {
		return fmt.Errorf("execute failed, err:%s", h.resp.GetErrInfo())
	}

	return nil
}
