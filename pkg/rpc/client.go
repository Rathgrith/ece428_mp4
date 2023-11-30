package rpc

import (
	"ece428_mp4/idl"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

const (
	DefaultRPCConnPort = 12315
)

type ClientManager struct {
	rpcConn   map[string]*grpc.ClientConn
	connMapMu sync.RWMutex
}

func NewRPCClientManager() *ClientManager {
	manager := ClientManager{
		rpcConn:   make(map[string]*grpc.ClientConn),
		connMapMu: sync.RWMutex{},
	}

	return &manager
}

// TODO: config
func getDefaultRPCConnAddress(host string) string {
	return fmt.Sprintf("%s:%d", host, DefaultRPCConnPort)
}

func (cm *ClientManager) newRPCConn(host string) (*grpc.ClientConn, error) {
	return grpc.Dial(getDefaultRPCConnAddress(host), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (cm *ClientManager) GetRPCConn(host string) (*grpc.ClientConn, error) {
	// TODO: cache conn handle
	return cm.newRPCConn(host)
}

func (cm *ClientManager) GetNodeManagerRPCClient(host string) (idl.NodeManageServiceClient, error) {
	return cm.newNodeManagerRPCClient(host)
}

func (cm *ClientManager) newNodeManagerRPCClient(host string) (idl.NodeManageServiceClient, error) {
	conn, err := cm.GetRPCConn(host)
	if err != nil {
		return nil, fmt.Errorf("can not dial to target address:%w", err)
	}

	client := idl.NewNodeManageServiceClient(conn)
	return client, nil
}
