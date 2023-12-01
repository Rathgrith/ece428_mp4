package juice

import (
	"bufio"
	"bytes"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	"ece428_mp4/pkg/maple_juice/juice/output"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type WorkerConfig struct {
	Func maple_juice.JuiceFunc
}

type TaskWorker struct {
	inputIntermediateFilenames []string
	fsClient                   *SDFSSDK.SDFSClient
	juiceFunc                  maple_juice.JuiceFunc
	attemptID                  string
	outputFilename             string

	localFileWriter *bufio.Writer
	fWriter         output.FormatWriter
}

func NewJuiceTaskWorker(inputIntermediateFilenames []string, client *SDFSSDK.SDFSClient, config *WorkerConfig,
	attemptID string, outputFilename string) (*TaskWorker, error) {

	localFile, err := os.Create("./" + outputFilename)
	if err != nil {
		return nil, fmt.Errorf("can not creat local output file:%w", err)
	}
	writer := bufio.NewWriter(localFile)

	worker := TaskWorker{
		inputIntermediateFilenames: inputIntermediateFilenames,
		fsClient:                   client,
		juiceFunc:                  config.Func,
		attemptID:                  attemptID,
		outputFilename:             outputFilename,
		fWriter:                    output.NewJsonMarshalFormatWriter(writer),
		localFileWriter:            writer,
	}

	return &worker, nil
}

func (tw *TaskWorker) Work() (*idl.RunJuiceTaskResponse, error) {
	for _, filename := range tw.inputIntermediateFilenames {
		err := tw.processSingleIntermediate(filename)
		if err != nil {
			return nil, fmt.Errorf("process input (%s) failed:%w", filename, err)
		}
	}

	err := tw.localFileWriter.Flush()
	if err != nil {
		return nil, fmt.Errorf("flush fail err:%w", err)
	}
	tmpName := tw.outputFilename + "-" + tw.attemptID
	err = tw.fsClient.TempPutLocalFile(tw.outputFilename, tw.attemptID, tw.outputFilename, "./")
	if err != nil {
		return nil, fmt.Errorf("upload temp output file failed:%w", err)
	}

	resp := idl.RunJuiceTaskResponse{
		Code:              idl.StatusCode_Success,
		TmpOutputFilename: tmpName,
	}

	return &resp, nil
}

func (tw *TaskWorker) processSingleIntermediate(filename string) error {
	// load all kv in current intermediate
	// TODO: consider large intermediate, currently just load in memory
	streamReader, err := tw.fsClient.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("can get stream reader for file (%s):%w", filename, err)
	}

	buf := bytes.Buffer{}

	for {
		content, err := streamReader.NextPacket()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read intermediate failed:%w", err)
		}

		if _, err = buf.Write(content); err != nil {
			return fmt.Errorf("can not write to buffer:%w", err)
		}
	}

	// extract kv from intermediate
	decoder := json.NewDecoder(&buf)
	kvs := make([]*maple_juice.KV, 0)
	for {
		kv := maple_juice.KV{}
		err := decoder.Decode(&kv)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("can not decode cur kv:%w", err)
		}

		kvs = append(kvs, &kv)
	}

	// do juice
	resKV, err := tw.juiceFunc(kvs)
	if err != nil {
		return fmt.Errorf("call juice func failed:%w", err)
	}

	// output to local file first
	err = tw.fWriter.Output(resKV)
	if err != nil {
		return fmt.Errorf("can not output current result kv:%w", err)
	}

	return nil
}
