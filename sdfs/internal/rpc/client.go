package rpc

import (
	idl2 "ece428_mp4/sdfs/internal/idl"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

const (
	DefaultRPCConnPort = 10086
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

func (m *ClientManager) newRPCConn(host string) (*grpc.ClientConn, error) {
	return grpc.Dial(getDefaultRPCConnAddress(host), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (m *ClientManager) GetRPCConn(host string) (*grpc.ClientConn, error) {
	// TODO: cache conn handle
	return m.newRPCConn(host)
}

func (m *ClientManager) GetDataNodePRCClient(host string) (idl2.DataNodeServiceClient, error) {
	return m.newDataNodeRPCClient(host)
}

func (m *ClientManager) GetNameNodePRCClient(host string) (idl2.NameNodeServiceClient, error) {
	return m.newNameNodeRPCClient(host)
}

func (m *ClientManager) newDataNodeRPCClient(host string) (idl2.DataNodeServiceClient, error) {
	conn, err := m.GetRPCConn(host)
	if err != nil {
		return nil, fmt.Errorf("can not dial to target address:%w", err)
	}

	client := idl2.NewDataNodeServiceClient(conn)
	return client, nil
}

func (m *ClientManager) newNameNodeRPCClient(host string) (idl2.NameNodeServiceClient, error) {
	conn, err := m.GetRPCConn(host)
	if err != nil {
		return nil, fmt.Errorf("can not dial to target address:%w", err)
	}

	client := idl2.NewNameNodeServiceClient(conn)
	return client, nil
}
