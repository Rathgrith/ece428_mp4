package main

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	"ece428_mp4/pkg/maple_juice/juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"

	"google.golang.org/protobuf/proto"
)

func main() {
	config := juice.WorkerConfig{
		Func: Juice,
	}
	RunJuice(&config)
}

func Juice(kvs []*maple_juice.KV) (*maple_juice.KV, error) {
	detectionCounts := make(map[string]int)
	total := 0

	// Aggregate counts for each Detection_ value
	for _, kv := range kvs {
		detectionValue := kv.Key
		count, _ := strconv.Atoi(kv.Value.(string))
		detectionCounts[detectionValue] += count
		total += count
	}

	// Calculate percentage composition
	var result string
	for detection, count := range detectionCounts {
		percentage := (float64(count) / float64(total)) * 100
		result += fmt.Sprintf("%s: %.2f%%\n", detection, percentage)
	}

	return &maple_juice.KV{Key: "Result", Value: result}, nil
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
