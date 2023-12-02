package main

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	"ece428_mp4/pkg/maple_juice/maple"
	"ece428_mp4/pkg/maple_juice/maple/input"
	"ece428_mp4/pkg/maple_juice/maple/partition"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/gob"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	config := maple.WorkerConfig{
		ReaderType:      input.LineRecordReaderType,
		Func:            Maple,
		PartitionerType: partition.PerKeyPartitionerType,
	}
	RunMaple(&config)
}

func Maple(kv *maple_juice.KV) (*maple_juice.KV, error) {
	// Implement your logic here, you can also modify the config
	panic("Implement me")
}

func RunMaple(config *maple.WorkerConfig) {
	decoder := gob.NewDecoder(os.Stdin)
	req := idl.RunMapleTaskRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		ReturnErrResponse(fmt.Errorf("can not read request:%w", err))
		return
	}

	worker, err := maple.NewTaskWorker(req.GetInputSplits(), SDFSSDK.NewSDFSClient(), config, req.GetAttemptId(),
		req.GetIntermediateFilenamePrefix())
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

func ReturnResponse(response *idl.RunMapleTaskResponse) {
	encoder := gob.NewEncoder(os.Stdout)
	err := encoder.Encode(response)
	if err != nil {
		return
	}
}

func ReturnErrResponse(err error) {
	ReturnResponse(&idl.RunMapleTaskResponse{
		Code:                 idl.StatusCode_InternalErr,
		TmpIntermediateFiles: nil,
		ErrInfo:              proto.String(err.Error()),
	})
}
