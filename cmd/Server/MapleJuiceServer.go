package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"

	"google.golang.org/grpc"

	"ece428_mp4/idl" // Adjust with the actual generated package path
)

// sample command: go run .\cmd\Server\MapleJuiceServer.go
type server struct {
	idl.UnimplementedMapleJuiceServiceServer
}

func (s *server) ExecuteMaple(ctx context.Context, req *idl.MapleRequest) (*idl.MapleResponse, error) {
	fmt.Printf("Received Maple request: %+v\n", req)
	if req.Exe == "joinMaple.exe" {
		fmt.Println("Join Maple")
		return &idl.MapleResponse{Output: "Join Maple"}, nil
	} else {
		fmt.Println("Filter Maple")
		output, err := runCommand(req.Exe, "-input", req.SrcDir1, "-prefix", req.Prefix, "-regex", req.Regex)
		if err != nil {
			return nil, err
		}
		return &idl.MapleResponse{Output: output}, nil
	}
}

func (s *server) ExecuteJuice(ctx context.Context, req *idl.JuiceRequest) (*idl.JuiceResponse, error) {
	fmt.Printf("Received Juice request: %+v\n", req)
	output, err := runCommand(req.Exe, "-key", req.Key, "-prefix", req.Prefix, "-output", req.OutDir, "-delete", req.DeleteInput)
	if err != nil {
		return nil, err
	}
	return &idl.JuiceResponse{Output: output}, nil
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	idl.RegisterMapleJuiceServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
