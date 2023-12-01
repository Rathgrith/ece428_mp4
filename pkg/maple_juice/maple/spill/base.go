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
	Commit() ([]string, error)
}

type LocalKVSpiller struct {
	partitionWriters map[string]*bufio.Writer
	encoders         map[string]*json.Encoder
	storeDir         string
	SDFCClient       *SDFSSDK.SDFSClient
	attemptID        string
}

func NewLocalKVSpiller(storeDir string, client *SDFSSDK.SDFSClient, ID string) *LocalKVSpiller {
	spiller := LocalKVSpiller{
		partitionWriters: make(map[string]*bufio.Writer),
		encoders:         make(map[string]*json.Encoder),
		storeDir:         storeDir,
		SDFCClient:       client,
		attemptID:        ID,
	}
	return &spiller
}

func (s *LocalKVSpiller) Spill(partitionID string, kv *maple_juice.KV) error {
	if _, exist := s.partitionWriters[partitionID]; !exist {
		localFileHandle, err := os.Create(s.storeDir + "/" + partitionID)
		if err != nil {
			return fmt.Errorf("can not create local spill file:%w", err)
		}
		writer := bufio.NewWriter(localFileHandle)
		s.partitionWriters[partitionID] = writer
		s.encoders[partitionID] = json.NewEncoder(writer)
	}

	err := s.encoders[partitionID].Encode(kv)
	if err != nil {
		return fmt.Errorf("can not spill current kv:%w", err)
	}

	return nil
}

func (s *LocalKVSpiller) Commit() ([]string, error) {
	partitionSDFSNames := make([]string, 0, len(s.partitionWriters))
	for partitionID, writer := range s.partitionWriters {
		err := writer.Flush()
		if err != nil {
			return nil, fmt.Errorf("can not write to %v : %w", partitionID, err)
		}
		err = s.SDFCClient.TempPutLocalFile(partitionID, s.attemptID, partitionID, s.storeDir)
		if err != nil {
			return nil, fmt.Errorf("can not upload partition (%s):%w", partitionID, err)
		}
		partitionSDFSNames = append(partitionSDFSNames, partitionID+"-"+s.attemptID)
	}
	return partitionSDFSNames, nil
}
