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
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

var col string
var dataset string

func main() {
	config := maple.WorkerConfig{
		ReaderType:      input.LineRecordReaderType,
		Func:            Maple,
		PartitionerType: partition.PerKeyPartitionerType,
	}
	col = os.Args[1]
	dataset = os.Args[2]
	// convert the string to int
	// specify the column to be used as the key
	RunMaple(&config)
}

func Maple(kv *maple_juice.KV) (*maple_juice.KV, error) {
	// v is the actual line of the csv file
	vList := strings.Split(kv.Value.(string), ",")
	idx, err := strconv.Atoi(col)
	if err != nil {
		return nil, err
	}
	joinKey := vList[idx]
	joinValue := dataset + "|" + kv.Value.(string)
	newKV := maple_juice.KV{
		Key:   joinKey,
		Value: joinValue,
	}
	return &newKV, nil
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
