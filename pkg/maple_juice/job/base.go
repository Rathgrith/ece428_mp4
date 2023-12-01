package job

import (
	"context"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/rpc"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"sync"
	"time"
)

const (
	timeout = time.Second * 15
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

func (m *Manager) SubmitMapleJob(req *idl.ExecuteMapleJobRequest) error {
	tracker := NewMapleJobTracker(req, m.rpcClientManager, m.fsClient, m)
	_, err := tracker.ExecuteJob()
	if err != nil {
		return err
	}
	return nil
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
	if t, exist := m.availableNodeHosts[host]; !exist {
		return false
	} else {
		return time.Now().Sub(t) < timeout
	}
}
