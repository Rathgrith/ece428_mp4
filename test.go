package main

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice/job"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"github.com/sirupsen/logrus"
)

func main() {
	logutil.InitDefaultLogger(logrus.DebugLevel)
	client := SDFSSDK.NewSDFSClient()

	inputFilename := "test.csv"
	mapleExe := "test_maple"

	err := client.PutLocalFile(mapleExe, mapleExe, "./workspace", true)
	if err != nil {
		panic(err)
	}
	err = client.PutLocalFile(inputFilename, inputFilename, "./workspace", true)
	if err != nil {
		panic(err)
	}

	jobManager := job.NewJobManager()
	jobManager.Heartbeat(context.Background(), &idl.HeartbeatRequest{Host: "fa23-cs425-4805.cs.illinois.edu"})
	mapleResp, err := jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
		ExeName:                    mapleExe,
		IntermediateFilenamePrefix: "TEST2",
		InputFiles:                 []string{inputFilename},
		NumMaples:                  3,
		ExeArgs:                    []string{"-regex test_regex"},
	})
	if err != nil || mapleResp.Code != idl.StatusCode_Success {
		panic(err)
	}
}
