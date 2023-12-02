package main

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/node"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"
)

func main() {
	StartRunNodeManager()
}

const (
	DefaultJobManagerHost = "fa23-cs425-4801.cs.illinois.edu:12301"
)

func StartRunNodeManager() {
	err := logutil.InitDefaultLogger(logrus.DebugLevel)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":12315")
	if err != nil {
		panic(err)
	}

	handler := node.NewManageHandler()
	err = handler.InitEnv()
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	idl.RegisterNodeManageServiceServer(server, handler)

	fmt.Println("node manager start to serve")

	go func() {
		selfHost, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		conn, err := grpc.Dial(DefaultJobManagerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		client := idl.NewJobManageServiceClient(conn)
		ticker := time.NewTicker(time.Second * 5)
		for {
			<-ticker.C
			client.Heartbeat(context.Background(), &idl.HeartbeatRequest{Host: selfHost})
		}
	}()

	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
