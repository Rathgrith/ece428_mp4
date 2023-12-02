package main

import (
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	logutil.InitDefaultLogger(logrus.DebugLevel)
	client := SDFSSDK.NewSDFSClient()

	//err := client.PutLocalFile("test_juice", "test_juice", "./workspace", true)
	//if err != nil {
	//	panic(err)
	//}
	//
	//nodeManager, err := rpc.NewRPCClientManager().GetNodeManagerRPCClient("fa23-cs425-4805.cs.illinois.edu")
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = nodeManager.RunJuiceTask(context.Background(), &idl.RunJuiceTaskRequest{
	//	ExeName:                "test_juice",
	//	InputIntermediateFiles: []string{"TEST2_odd", "TEST2_even"},
	//	OutPutFilename:         "TEST_JUICE",
	//	AttemptId:              "JUICE1_1",
	//})
	//if err != nil {
	//	panic(err)
	//}

	inputFilename := "test.csv"
	// mapleExe := "filterMaple.exe"
	// juiceExe := "filterJuice.exe"

	// err := client.PutLocalFile(mapleExe, mapleExe, "./workspace", true)
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.PutLocalFile(juiceExe, juiceExe, "./workspace", true)
	// if err != nil {
	// 	panic(err)
	// }

	err := client.PutLocalFile(inputFilename, inputFilename, "./workspace", true)
	if err != nil {
		panic(err)
	}

	// jobManager := job.NewJobManager()
	// jobManager.Heartbeat(context.Background(), &idl.HeartbeatRequest{Host: "fa23-cs425-4805.cs.illinois.edu"})
	// mapleResp, err := jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
	// 	ExeName:                    mapleExe,
	// 	IntermediateFilenamePrefix: "TEST2",
	// 	InputFiles:                 []string{inputFilename},
	// 	NumMaples:                  3,
	// 	ExeArgs:                    []string{"-regex test_regex"},
	// })
	// if err != nil || mapleResp.Code != idl.StatusCode_Success {
	// 	panic(err)
	// }

	// juiceResp, err := jobManager.SubmitJuiceJob(&idl.ExecuteJuiceJobRequest{
	// 	ExeName:               juiceExe,
	// 	IntermediateFilenames: mapleResp.GetIntermediateFilenames(),
	// 	NumMaples:             2,
	// 	OutPutFilename:        "TEST_JUICE",
	// 	ExeArgs:               nil,
	// })
	// if err != nil || juiceResp.Code != idl.StatusCode_Success {
	// 	panic(err)
	// }

}

func testJuice(kvs []*maple_juice.KV) (*maple_juice.KV, error) {
	count := len(kvs)

	newKV := maple_juice.KV{
		Key:   "Count",
		Value: strconv.Itoa(count),
	}

	return &newKV, nil
}
