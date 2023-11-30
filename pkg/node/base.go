package node

import (
	"ece428_mp4/idl"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"fmt"
	"os"
)

const (
	DefaultPerm = 0755
)

const (
	DefaultStoreDir    = "./workspace"
	DefaultExeStoreDir = DefaultStoreDir + "/exe"
)

type ManageServiceHandler struct {
	idl.NodeManageServiceServer
	fsClient *SDFSSDK.SDFSClient
}

func NewManageHandler() *ManageServiceHandler {
	handler := ManageServiceHandler{
		fsClient: SDFSSDK.NewSDFSClient(),
	}

	return &handler
}

func (h *ManageServiceHandler) InitEnv() error {
	for _, dir := range []string{
		DefaultStoreDir, DefaultExeStoreDir,
	} {
		err := os.MkdirAll(dir, DefaultPerm)
		if err != nil {
			return fmt.Errorf("creat dir (%s) failed:%w", dir, err)
		}
	}

	return nil
}
