package main

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	"ece428_mp4/pkg/maple_juice/juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

func main() {
	config := juice.WorkerConfig{
		Func: Juice,
	}
	RunJuice(&config)
}

func Juice(kvs []*maple_juice.KV) (*maple_juice.KV, error) {
	// for all values split by | ,
	// if the first part is D1, add it to list D1
	// if the first part is D2, add it to list D2
	D1 := make([]string, 0)
	D2 := make([]string, 0)
	for _, kv := range kvs {
		value := kv.Value.(string)
		if strings.HasPrefix(value, "D1") {
			D1 = append(D1, value)
		} else if strings.HasPrefix(value, "D2") {
			D2 = append(D2, value)
		}
	}
	// while D1 and D2 are not empty
	// pop a value from D1 and D2
	// concatenate them with | and output
	var output string
	if len(D1) == 0 || len(D2) == 0 {
		return nil, nil
	}
	lineCount := len(D1)
	if len(D2) < lineCount {
		lineCount = len(D2)
	}
	for i := 0; i < lineCount; i++ {
		output += D1[i] + "|" + D2[i] + "\n"
	}
	newKV := maple_juice.KV{
		Key:   "",
		Value: output,
	}
	return &newKV, nil
}
func RunJuice(config *juice.WorkerConfig) {
	decoder := gob.NewDecoder(os.Stdin)
	req := idl.RunJuiceTaskRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		ReturnErrResponse(fmt.Errorf("can not read request:%w", err))
		return
	}

	worker, err := juice.NewJuiceTaskWorker(req.GetInputIntermediateFiles(), SDFSSDK.NewSDFSClient(), config,
		req.GetAttemptId(), req.GetOutPutFilename())
	if err != nil {
		ReturnErrResponse(fmt.Errorf("can init task worker:%w", err))
		return
	}

	resp, err := worker.Work()
	if err != nil {
		ReturnErrResponse(err)
		return
	}
	ReturnResponse(resp)
}

func ReturnResponse(response *idl.RunJuiceTaskResponse) {
	encoder := gob.NewEncoder(os.Stdout)
	err := encoder.Encode(response)
	if err != nil {
		return
	}
}

func ReturnErrResponse(err error) {
	ReturnResponse(&idl.RunJuiceTaskResponse{
		Code:              idl.StatusCode_InternalErr,
		TmpOutputFilename: "",
		ErrInfo:           proto.String(err.Error()),
	})
}
