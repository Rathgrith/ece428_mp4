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
	exeFile := "test_exe"

	err := client.PutLocalFile(exeFile, exeFile, "./workspace", true)
	if err != nil {
		panic(err)
	}

	err = client.PutLocalFile(inputFilename, inputFilename, "./workspace", true)
	if err != nil {
		panic(err)
	}

	jobManager := job.NewJobManager()
	jobManager.Heartbeat(context.Background(), &idl.HeartbeatRequest{Host: "fa23-cs425-4805.cs.illinois.edu"})
	err = jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
		ExeName:                    exeFile,
		IntermediateFilenamePrefix: "TEST2",
		InputFiles:                 []string{inputFilename},
		NumMaples:                  3,
	})
	if err != nil {
		panic(err)
	}

	return

	//inputFileSize, err := client.GetFileSize(inputFilename)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.PutLocalFile(exeFile, exeFile, "./workspace", true)
	//if err != nil {
	//	panic(err)
	//}
	//
	//rpcClientManager := rpc.NewRPCClientManager()
	//nodeManagerClient, err := rpcClientManager.GetNodeManagerRPCClient("fa23-cs425-4805.cs.illinois.edu")
	//if err != nil {
	//	panic(err)
	//}
	//
	//req := idl.RunMapleTaskRequest{
	//	ExeName:                    exeFile,
	//	IntermediateFilenamePrefix: "hello",
	//	InputSplits: []*idl.InputSplit{{
	//		InputFileName: "test.csv",
	//		InputOffset:   0,
	//		InputSize:     int32(inputFileSize),
	//	}},
	//	AttemptId: "test",
	//}
	//
	//resp, err := nodeManagerClient.RunMapleTask(context.Background(), &req)
	//if err != nil {
	//	panic(err)
	//}
	//if resp.GetCode() != idl.StatusCode_Success {
	//	panic(errors.New(resp.GetErrInfo()))
	//}
	//
	//fmt.Println(resp.GetTmpIntermediateFiles())

	//err = client.PutLocalFile(inputFilename, inputFilename, "./workspace", true)
	//if err != nil {
	//	panic(err)
	//}
	//
	//exeFile := "test_exe"
	//err = client.PutLocalFile(exeFile, exeFile, "./workspace", true)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.GetFileToLocal(exeFile, exeFile, "./workspace/exe")
	//if err != nil {
	//	panic(err)
	//}

	//s := split.NewRawFileSplitter([]string{inputFilename}, client, 4)
	//splits, err := s.Split()
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, cur := range splits {
	//	b, err := json.Marshal(cur.Split)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(string(b))
	//	recordReader, err := input.NewLineRecordReader(client, cur.Split)
	//	if err != nil {
	//		panic(err)
	//	}
	//	for {
	//		kv, err := recordReader.NextKeyValue()
	//		if err == io.EOF {
	//			break
	//		}
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Printf("Key:%s, Value:%s\n", kv.Key, string(kv.Value.([]byte)))
	//	}
	//	fmt.Println(cur.Locations)
	//}

	//
	//resp, err := node.NewRunMapleTaskHandler(context.Background(), &req, client).Handle()
	//if err != nil {
	//	panic(err)
	//}
	//
	//marshal, err := json.Marshal(resp)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(marshal)

	//err := client.TouchFile("1.csv")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.TempPutLocalFile("1.csv", "worker1", "1.csv", "./workspace")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.TempPutLocalFile("1.csv", "worker2", "2.csv", "./workspace")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.TempPutLocalFile("1.csv", "worker3", "3.csv", "./workspace")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.MergeFiles("1.csv", []string{"1.csv-worker1", "1.csv-worker2", "1.csv-worker3"}, true, true)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.ReadFileToLocal("1.csv", 0, 0, "test.csv", "./workspace")
	//if err != nil {
	//	panic(err)
	//}

	//err = client.PutLocalFile(inputFilename, inputFilename, "./workspace")
	//
	//_, err := client.GetFileSize(inputFilename)
	//if err != nil {
	//	panic(err)
	//}
}
