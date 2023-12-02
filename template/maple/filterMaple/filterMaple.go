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
	"regexp"

	"google.golang.org/protobuf/proto"
)

var regexCondition string

func main() {
	config := maple.WorkerConfig{
		ReaderType:      input.LineRecordReaderType,
		Func:            Maple,
		PartitionerType: partition.PerKeyPartitionerType,
	}
	regexCondition = os.Args[1]
	RunMaple(&config)
}

func Maple(kv *maple_juice.KV) (*maple_juice.KV, error) {
	// Implement your logic here, you can also modify the config
	val, valid := kv.Value.([]byte)
	if !valid {
		return nil, fmt.Errorf("can not convert value")
	}
	v := string(val)
	var newKey string
	// if regexCondition can match v then return v as new key
	if matched, _ := regexp.MatchString(regexCondition, v); matched {
		newKey = v
	} else {
		return nil, nil
	}
	// else skip this key-value pair

	newKV := maple_juice.KV{
		Key:   newKey,
		Value: string(kv.Value.([]byte)),
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
