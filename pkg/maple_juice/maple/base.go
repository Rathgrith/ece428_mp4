package maple

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	"ece428_mp4/pkg/maple_juice/maple/input"
	"ece428_mp4/pkg/maple_juice/maple/partition"
	"ece428_mp4/pkg/maple_juice/maple/spill"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"fmt"
	"io"
)

type WorkerConfig struct {
	ReaderType      input.RecordReaderType
	Func            maple_juice.MapleFunc
	PartitionerType partition.PartitionerType
}

type TaskWorker struct {
	fsClient            *SDFSSDK.SDFSClient
	recordReaderNewFunc input.NewRecordReaderFunc
	mapleFunc           maple_juice.MapleFunc
	partitioner         partition.Partitioner
	spiller             spill.KVSpiller

	inputSplits        []*idl.InputSplit
	attemptID          string
	intermediatePrefix string
}

func NewTaskWorker(inputSplits []*idl.InputSplit, client *SDFSSDK.SDFSClient, config *WorkerConfig,
	attemptID string, intermediatePrefix string) (*TaskWorker, error) {
	readerNewFunc := input.GetNewRecordReaderFunc(config.ReaderType)
	if readerNewFunc == nil {
		return nil, fmt.Errorf("do not specifiy record reader type")
	}

	partitionerNewFunc := partition.GetNewRecordReaderFunc(config.PartitionerType)
	if readerNewFunc == nil {
		return nil, fmt.Errorf("do not specifiy partitioner type")
	}
	partitioner := partitionerNewFunc()

	worker := TaskWorker{
		inputSplits:         inputSplits,
		fsClient:            client,
		recordReaderNewFunc: readerNewFunc,
		mapleFunc:           config.Func,
		partitioner:         partitioner,
		spiller:             spill.NewLocalKVSpiller("./", client, attemptID),
		attemptID:           attemptID,
		intermediatePrefix:  intermediatePrefix,
	}

	return &worker, nil
}

func (tw *TaskWorker) Work() (*idl.RunMapleTaskResponse, error) {
	for _, inputSplit := range tw.inputSplits {
		if err := tw.processSingleInputSplit(inputSplit); err != nil {
			return nil, fmt.Errorf("raise err during maple:%w", err)
		}
	}

	intermediates, err := tw.spiller.Commit()
	if err != nil {
		return nil, fmt.Errorf("submit spill file failed:%w", err)
	}

	resp := idl.RunMapleTaskResponse{
		Code:                 idl.StatusCode_Success,
		TmpIntermediateFiles: intermediates,
		ErrInfo:              nil,
	}

	return &resp, nil
}

func (tw *TaskWorker) processSingleInputSplit(split *idl.InputSplit) error {
	recordReader, err := tw.recordReaderNewFunc(tw.fsClient, split)
	if err != nil {
		return fmt.Errorf("can not get read stream for file (%s):%w", split.InputFileName, err)
	}

	mapleFunc := tw.mapleFunc
	partitioner := tw.partitioner
	spiller := tw.spiller

	for {
		kv, err := recordReader.NextKeyValue()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("can not read next key-value:%w", err)
		}

		// do maple
		newKV, err := mapleFunc(kv)
		if err != nil {
			return fmt.Errorf("can not maple current key-value:%w", err)
		}

		// partition and spill
		if newKV != nil {
			partitionID := partitioner.GetPartitionName(newKV)
			err := spiller.Spill(tw.intermediatePrefix+"_"+partitionID, newKV)
			if err != nil {
				return fmt.Errorf("can not spill current key-value:%w", err)
			}
		}
	}
}
