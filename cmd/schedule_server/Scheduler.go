package main

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/maple_juice/job"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type server struct {
	idl.UnimplementedMapleJuiceSchedulerServer
	taskQueue  chan Task // Channel used as a FIFO queue
	jobManager *job.Manager
}

func (s *server) EnqueueTask(ctx context.Context, req *idl.TaskRequest) (*idl.TaskResponse, error) {
	fmt.Printf("Received task request: %+v\n", req)
	fmt.Printf("Task queue length: %d\n", len(s.taskQueue))
	// Create a channel for task completion signal
	completion := make(chan string, 1)
	task := Task{
		Type:          req.TaskType,
		Executable:    req.Exe,
		NumTasks:      int(req.NumJobs),
		Prefix:        req.Prefix,
		SrcDir1:       req.SrcDir1,
		SrcDir2:       req.SrcDir2,
		Regex:         req.Regex,
		JoinColumn1:   req.JoinColumn1,
		JoinColumn2:   req.JoinColumn2,
		OutDir:        req.DestFile,
		completionSig: completion,
	}
	fmt.Printf("Task: %+v\n", task)
	// Enqueue the task
	s.taskQueue <- task
	result := <-completion
	close(completion)

	return &idl.TaskResponse{Message: result}, nil
	// return &idl.TaskResponse{Message: "Task enqueued successfully"}, nil
}

// Task represents a MapleJuice task
type Task struct {
	Type          string // "maple" or "juice"
	Executable    string
	NumTasks      int
	Prefix        string
	SrcDir1       string
	SrcDir2       string
	Regex         string
	JoinColumn1   int32
	JoinColumn2   int32
	OutDir        string
	completionSig chan string // Channel to signal completion
}

// executeTask simulates task execution
func executeTask(jobManager *job.Manager, task Task) {
	fmt.Printf("Executing task: %+v\n", task)
	if task.Executable == "filterMaple.exe" {
		err := jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
			ExeName:                    task.Executable,
			IntermediateFilenamePrefix: task.Prefix,
			InputFiles:                 []string{task.SrcDir1},
			NumMaples:                  int32(task.NumTasks),
			ExeArgs:                    []string{task.Regex},
		})
		if err != nil {
			panic(err)
		}
		// } else if task.Executable == "filterJuice.exe" {
		// 	err := jobManager.SubmitJuiceJob(&idl.ExecuteJuiceJobRequest{
		// 		ExeName:                task.Executable,
		// 		InputIntermediateFiles: []string{task.Prefix},
		// 		OutPutFilename:         task.OutDir,
		// 		AttemptId:              "JUICE1_1",
		// 	})
		// 	if err != nil {
		// 		panic(err)
		// 	}
	} else if task.Executable == "joinMaple.exe" {
		err := jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
			ExeName:                    task.Executable,
			IntermediateFilenamePrefix: task.Prefix,
			InputFiles:                 []string{task.SrcDir1},
			NumMaples:                  int32(task.NumTasks),
			ExeArgs:                    []string{strconv.Itoa(int(task.JoinColumn1)), " D1"},
		})
		if err != nil {
			panic(err)
		}
		err = jobManager.SubmitMapleJob(&idl.ExecuteMapleJobRequest{
			ExeName:                    task.Executable,
			IntermediateFilenamePrefix: task.Prefix,
			InputFiles:                 []string{task.SrcDir2},
			NumMaples:                  int32(task.NumTasks),
			ExeArgs:                    []string{strconv.Itoa(int(task.JoinColumn2)), " D2"},
		})
		if err != nil {
			panic(err)
		}
		// } else if task.Executable == "joinJuice.exe" {
		// 	err := jobManager.SubmitJuiceJob(&idl.ExecuteJuiceJobRequest{
		// 		ExeName:                task.Executable,
		// 		InputIntermediateFiles: []string{task.Prefix},
		// 		OutPutFilename:         task.OutDir,
		// 		AttemptId:              "JUICE1_1",
		// 	})
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }
	}
	fmt.Printf("Task completed: %+v\n", task)
	// Signal task completion
	task.completionSig <- "Task completed successfully with output file: " + task.OutDir
}

func main() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}
	taskQueue := make(chan Task, 100) // Task queue with a buffer of 100 tasks

	// Create and start the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sqlServer := &server{taskQueue: taskQueue, jobManager: job.NewJobManager()}
	idl.RegisterMapleJuiceSchedulerServer(grpcServer, sqlServer)

	// Start the scheduler in a separate goroutine
	go func() {
		for task := range taskQueue {
			executeTask(sqlServer.jobManager, task)
		}
	}()

	log.Println("Starting gRPC server on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
