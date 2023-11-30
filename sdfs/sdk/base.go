package sdk

import (
	"context"
	idl2 "ece428_mp4/sdfs/internal/idl"
	"ece428_mp4/sdfs/internal/rpc"
	"fmt"
)

const (
	DefaultNameNodeHost = "fa23-cs425-4801.cs.illinois.edu"
)

type SDFSClient struct {
	rpcManager *rpc.ClientManager
}

func NewSDFSClient() *SDFSClient {
	client := SDFSClient{
		rpcManager: rpc.NewRPCClientManager(),
	}

	return &client
}

func (c *SDFSClient) getNameNodeClient() (idl2.NameNodeServiceClient, error) {
	// TODO: ask name node host instead hard code
	return c.rpcManager.GetNameNodePRCClient(DefaultNameNodeHost)
}

func (c *SDFSClient) GetFileLocationHosts(filename string) ([]string, error) {
	locResp, err := c.getFileLocations(filename)
	if err != nil {
		return nil, err
	}

	locHosts := make([]string, 0)
	locHosts = append(locHosts, locResp.GetCoordinator().GetHostname())
	for _, node := range locResp.GetSecondaries() {
		locHosts = append(locHosts, node.GetHostname())
	}

	return locHosts, nil
}

func (c *SDFSClient) getFileLocations(filename string) (*idl2.QueryFileResponse, error) {
	nameNode, err := c.getNameNodeClient()
	if err != nil {
		return nil, fmt.Errorf("can not dial name node:%w", err)
	}

	queryReq := idl2.QueryFileRequest{Filename: filename}
	return nameNode.QueryFile(context.Background(), &queryReq)
}

func (c *SDFSClient) askFileStoreLocations(filename string) (*idl2.AskFileStoreResponse, error) {
	nameNode, err := c.getNameNodeClient()
	if err != nil {
		return nil, fmt.Errorf("can not dial name node:%w", err)
	}

	queryReq := idl2.AskFileStoreRequest{Filename: filename}
	return nameNode.AskFileStore(context.Background(), &queryReq)
}
