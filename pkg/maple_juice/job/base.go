package job

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/rpc"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type TaskStatus int

const (
	Init TaskStatus = iota + 1
	Running
	Errored
	Failed
)

const (
	ServeHost = 12301
)

const (
	timeout      = time.Second * 15
	MaxRetryTime = 10
)

type Manager struct {
	idl.JobManageServiceServer

	rpcClientManager   *rpc.ClientManager
	fsClient           *SDFSSDK.SDFSClient
	availableNodeHosts map[string]time.Time
	mu                 sync.RWMutex
}

func NewJobManager() *Manager {
	manager := Manager{
		rpcClientManager:   rpc.NewRPCClientManager(),
		fsClient:           SDFSSDK.NewSDFSClient(),
		availableNodeHosts: make(map[string]time.Time),
		mu:                 sync.RWMutex{},
	}

	return &manager
}

func (m *Manager) StartServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", ServeHost))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	idl.RegisterJobManageServiceServer(server, m)

	fmt.Println("job Manager start to serve")

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (m *Manager) SubmitMapleJob(req *idl.ExecuteMapleJobRequest) (*idl.ExecuteMapleJobResponse, error) {
	tracker := NewMapleJobTracker(context.Background(), req, m.rpcClientManager, m.fsClient, m)
	resp, err := tracker.ExecuteJob()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (m *Manager) SubmitJuiceJob(req *idl.ExecuteJuiceJobRequest) (*idl.ExecuteJuiceJobResponse, error) {
	tracker := NewJuiceJobTracker(context.Background(), req, m.rpcClientManager, m.fsClient, m)
	resp, err := tracker.ExecuteJob()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (m *Manager) Heartbeat(ctx context.Context, request *idl.HeartbeatRequest) (*idl.HeartBeatResponse, error) {
	m.mu.Lock()
	m.availableNodeHosts[request.GetHost()] = time.Now()
	m.mu.Unlock()
	return &idl.HeartBeatResponse{}, nil
}

func (m *Manager) GetAvailableHost() []string {
	hosts := make([]string, 0)
	curTime := time.Now()
	m.mu.RLock()
	for host, t := range m.availableNodeHosts {
		if curTime.Sub(t) > timeout {
			continue
		}
		hosts = append(hosts, host)
	}
	m.mu.RUnlock()
	return hosts
}

func (m *Manager) NodeAlive(host string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if t, exist := m.availableNodeHosts[host]; !exist {
		return false
	} else {
		return time.Now().Sub(t) < timeout
	}
}
