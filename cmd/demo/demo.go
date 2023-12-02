package main

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice/job"
	"flag"

	"github.com/sirupsen/logrus"
)

func main() {
	interconne_type := flag.String("inter", "None", "interconnection type")
	flag.Parse()
	logutil.InitDefaultLogger(logrus.DebugLevel)
	// client := SDFSSDK.NewSDFSClient()

	// //err := client.PutLocalFile("test_juice", "test_juice", "./workspace", true)
	// //if err != nil {
	// //	panic(err)
	// //}
	// //
	// //nodeManager, err := rpc.NewRPCClientManager().GetNodeManagerRPCClient("fa23-cs425-4805.cs.illinois.edu")
	// //if err != nil {
	// //	panic(err)
	// //}
	// //
	// //_, err = nodeManager.RunJuiceTask(context.Background(), &idl.RunJuiceTaskRequest{
	// //	ExeName:                "test_juice",
	// //	InputIntermediateFiles: []string{"TEST2_odd", "TEST2_even"},
	// //	OutPutFilename:         "TEST_JUICE",
	// //	AttemptId:              "JUICE1_1",
	// //})
	// //if err != nil {
	// //	panic(err)
	// //}

	// inputFilename := "test.csv"
	// inputFilename2 := "Traffic_Signal_Intersections.csv"
	// mapleExe := "filterMaple"
	// juiceExe := "filterJuice"
	// mapleExe1 := "demoMaple"
	// juiceExe1 := "demoJuice"

	// err := client.PutLocalFile(mapleExe, mapleExe, "./", true)
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.PutLocalFile(juiceExe, juiceExe, "./", true)
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.PutLocalFile(inputFilename, inputFilename, "./", true)
	// if err != nil {
	// 	panic(err)
	// }
	// err = client.PutLocalFile(inputFilename2, inputFilename2, "./", true)
	// if err != nil {
	// 	panic(err)
	// }
	// err = client.PutLocalFile(mapleExe1, mapleExe1, "./", true)
	// if err != nil {
	// 	panic(err)
	// }
	// err = client.PutLocalFile(juiceExe1, juiceExe1, "./", true)
	// if err != nil {
	// 	panic(err)
	// }
	// client.GetFileToLocal("output.csv", "output.csv", "./")
	inter := *interconne_type

	jobManager := job.NewJobManager()
	jobManager.Heartbeat(context.Background(), &idl.HeartbeatRequest{Host: "fa23-cs425-4805.cs.illinois.edu"})
	mapleResp, err := jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
		ExeName:                    "demoMaple",
		IntermediateFilenamePrefix: "demo",
		InputFiles:                 []string{"Traffic_Signal_Intersections.csv"},
		NumMaples:                  3,
		ExeArgs:                    []string{inter},
	})
	if err != nil || mapleResp.Code != idl.StatusCode_Success {
		panic(err)
	}
	if len(mapleResp.GetIntermediateFilenames()) == 0 {
		logutil.Logger.Debugf("no lines match the regex")
	}
	juiceResp, err := jobManager.SubmitJuiceJob(&idl.ExecuteJuiceJobRequest{
		ExeName:               "demoJuice",
		IntermediateFilenames: mapleResp.GetIntermediateFilenames(),
		NumMaples:             4,
		OutPutFilename:        "output.csv",
		ExeArgs:               nil,
	})
	if err != nil || juiceResp.Code != idl.StatusCode_Success {
		logutil.Logger.Debugf("err:%v", err)
	}

}

// func testJuice(kvs []*maple_juice.KV) (*maple_juice.KV, error) {
// 	count := len(kvs)

// 	newKV := maple_juice.KV{
// 		Key:   "Count",
// 		Value: strconv.Itoa(count),
// 	}

// 	return &newKV, nil
// }
