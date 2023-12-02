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
	"strings"

	"google.golang.org/protobuf/proto"
)

var interconneType string

func main() {
	config := maple.WorkerConfig{
		ReaderType:      input.LineRecordReaderType,
		Func:            Maple,
		PartitionerType: partition.PerKeyPartitionerType,
	}
	interconneType = os.Args[1]
	RunMaple(&config)
}

func Maple(kv *maple_juice.KV) (*maple_juice.KV, error) {
	val, valid := kv.Value.([]byte)
	if !valid {
		return nil, fmt.Errorf("can not convert value")
	}
	v := string(val)
	// Split the CSV line into columns
	columns := strings.Split(v, ",")
	interconneCol := columns[10]
	detectionCol := columns[9]

	if interconneCol == interconneType {
		return &maple_juice.KV{Key: detectionCol, Value: "1"}, nil
	}

	return nil, nil // Skip this line if it doesn't match the Interconne type
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
