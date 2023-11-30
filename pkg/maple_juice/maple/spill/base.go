package spill

import (
	"bufio"
	"ece428_mp4/pkg/maple_juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/json"
	"fmt"
	"os"
)

type KVSpiller interface {
	Spill(string, *maple_juice.KV) error
	Submit() ([]string, error)
}

type LocalKVSpiller struct {
	partitions map[string]*bufio.Writer
	storeDir   string
	SDFCClient *SDFSSDK.SDFSClient
	taskID     string
}

func NewLocalKVSpiller(storeDir string, client *SDFSSDK.SDFSClient, ID string) *LocalKVSpiller {
	spiller := LocalKVSpiller{
		partitions: make(map[string]*bufio.Writer),
		storeDir:   storeDir,
		SDFCClient: client,
		taskID:     ID,
	}
	return &spiller
}

func (s *LocalKVSpiller) Spill(partitionID string, kv *maple_juice.KV) error {
	os.MkdirAll(s.storeDir, 0777)
	if _, exist := s.partitions[partitionID]; !exist {
		localFileHandle, err := os.Create(s.storeDir + "/" + partitionID)
		if err != nil {
			return fmt.Errorf("can not create local spill file:%w", err)
		}
		s.partitions[partitionID] = bufio.NewWriter(localFileHandle)
	}

	content, err := json.Marshal(kv)
	if err != nil {
		return fmt.Errorf("can not marshal kv:%w", err)
	}

	if _, err = s.partitions[partitionID].Write(content); err != nil {
		return fmt.Errorf("can not spill current kv:%w", err)
	}

	return nil
}

func (s *LocalKVSpiller) Submit() ([]string, error) {
	partitionSDFSNames := make([]string, 0, len(s.partitions))
	for partitionID, writer := range s.partitions {
		err := writer.Flush()
		if err != nil {
			return nil, fmt.Errorf("can not write to %v : %w", partitionID, err)
		}
		err = s.SDFCClient.TempPutLocalFile(partitionID, s.taskID, partitionID, s.storeDir)
		if err != nil {
			return nil, fmt.Errorf("can not upload partition (%s):%w", partitionID, err)
		}
		partitionSDFSNames = append(partitionSDFSNames, partitionID+"-"+s.taskID)
	}
	return partitionSDFSNames, nil
}
