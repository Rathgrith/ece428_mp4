package main

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/logutil"
	"ece428_mp4/pkg/node"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	StartRunNodeManager()
}

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

	fmt.Println("start to serve")

	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
