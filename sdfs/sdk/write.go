package sdk

import (
	"bytes"
	"context"
	idl2 "ece428_mp4/sdfs/internal/idl"
	"fmt"
	"io"
	"os"
)

const (
	DefaultPacketSize = 1024 * 64 // byte
)

type StreamWriter struct {
	stream idl2.DataNodeService_StoreFileClient
	buf    bytes.Buffer
}

func NewStreamWriter(stream idl2.DataNodeService_StoreFileClient) *StreamWriter {
	writer := StreamWriter{
		stream: stream,
		buf:    bytes.Buffer{},
	}

	return &writer
}

func (w *StreamWriter) Write(content []byte) error {
	// TODO: buffer
	if err := w.stream.Send(&idl2.StoreFileRequest{
		Content: content,
	}); err != nil {
		return err
	}
	return nil
}

func (w *StreamWriter) Close() error {
	resp, err := w.stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if resp.GetCode() != idl2.StatusCode_Success {

	}
	return nil
}

func (c *SDFSClient) GetFileWriteStream(filename string) (*StreamWriter, error) {
	storeLocResp, err := c.askFileStoreLocations(filename)
	if err != nil {
		return nil, fmt.Errorf("can not allocate file store node:%w", err)
	}

	coordinatorHost := storeLocResp.Coordinator.GetHostname()
	coordinator, err := c.rpcManager.GetDataNodePRCClient(coordinatorHost)
	if err != nil {
		return nil, fmt.Errorf("can not dial data node:%w", err)
	}

	stream, err := coordinator.StoreFile(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can not request store coordinator:%w", err)
	}

	metaRequest := idl2.StoreFileRequest{
		Filename:     filename,
		ReplicaNodes: storeLocResp.GetSecondaries(),
		Content:      nil,
	}
	err = stream.Send(&metaRequest)
	if err != nil {
		return nil, fmt.Errorf("can not init stream from coordinator:%w", err)
	}

	return NewStreamWriter(stream), nil
}

func (c *SDFSClient) PutLocalFile(filename string, localFilename string, storeDir string, override bool) error {
	if override {
		err := c.DeleteFile(filename)
		if err != nil {
			return fmt.Errorf("clear previous file failed:%w", err)
		}
	}
	writer, err := c.GetFileWriteStream(filename)
	if err != nil {
		return fmt.Errorf("can not get write stream:%w", err)
	}

	f, err := os.Open(storeDir + "/" + localFilename)
	if err != nil {
		return fmt.Errorf("can not open file:%w", err)
	}

	sendBuf := make([]byte, DefaultPacketSize)

	for {
		l, err := f.Read(sendBuf)
		if err == io.EOF {
			err = writer.Close()
			if err != nil {
				return fmt.Errorf("close stream failed:%w", err)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("read from local file failed:%w", err)
		}

		if err = writer.Write(sendBuf[:l]); err != nil {
			return fmt.Errorf("write file failed:%w", err)
		}
	}

	return nil
}

func (c *SDFSClient) DeleteFile(filename string) error {
	nameNode, err := c.getNameNodeClient()
	if err != nil {
		return fmt.Errorf("can not dial name node:%w", err)
	}

	req := idl2.DeleteFileRequest{Filename: filename}
	_, err = nameNode.DeleteFile(context.Background(), &req)
	if err != nil {
		// TODO: use status code
		if err.Error() != "rpc error: code = Unknown desc = file does not exist" {
			return fmt.Errorf("delete failed:%w", err)
		}
	}
	//if resp.GetCode() != idl.StatusCode_Success {
	//	return fmt.Errorf("delete unsuccess")
	//}
	return nil
}

func (c *SDFSClient) MergeFiles(outFilename string, mergeFiles []string, delete bool, padNewLineSymbol bool) error {
	// TODO: Ask namenode to merge (determine merge location)
	filename := mergeFiles[0]

	locResp, err := c.getFileLocations(filename)
	if err != nil {
		return fmt.Errorf("query file location failed:%w", err)
	}

	coordinatorHost := locResp.Coordinator.GetHostname()
	coordinator, err := c.rpcManager.GetDataNodePRCClient(coordinatorHost)
	if err != nil {
		return fmt.Errorf("can not dial data node:%w", err)
	}

	mergeReq := idl2.MergeFilesRequest{
		OutFilename:      outFilename,
		FilenameSeq:      mergeFiles,
		ReplicaNodes:     locResp.GetSecondaries(),
		Delete:           delete,
		PadNewLineSymbol: padNewLineSymbol,
	}
	resp, err := coordinator.MergeFiles(context.Background(), &mergeReq)
	if err != nil || resp.Code != idl2.StatusCode_Success {
		return fmt.Errorf("merge failed:%w", err)
	}

	return nil
}

func (c *SDFSClient) TouchFile(filename string) error {
	writer, err := c.GetFileWriteStream(filename)
	if err != nil {
		return fmt.Errorf("can not get write stream:%w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("close stream failed:%w", err)
	}
	return nil
}

func (c *SDFSClient) GetTempFileWriteStream(filename string, tempSuffix string) (*StreamWriter, error) {
	storeLocResp, err := c.askFileStoreLocations(filename)
	if err != nil {
		return nil, fmt.Errorf("can not allocate file store node:%w", err)
	}

	coordinatorHost := storeLocResp.Coordinator.GetHostname()
	coordinator, err := c.rpcManager.GetDataNodePRCClient(coordinatorHost)
	if err != nil {
		return nil, fmt.Errorf("can not dial data node:%w", err)
	}

	stream, err := coordinator.StoreFile(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can not request store coordinator:%w", err)
	}

	tempName := filename + "-" + tempSuffix
	metaRequest := idl2.StoreFileRequest{
		Filename:     tempName,
		ReplicaNodes: storeLocResp.GetSecondaries(),
		Content:      nil,
	}
	err = stream.Send(&metaRequest)
	if err != nil {
		return nil, fmt.Errorf("can not init stream from coordinator:%w", err)
	}

	return NewStreamWriter(stream), nil
}

func (c *SDFSClient) TempPutLocalFile(filename string, tempSuffix string, localFilename string, storeDir string) error {
	writer, err := c.GetTempFileWriteStream(filename, tempSuffix)
	if err != nil {
		return fmt.Errorf("can not get write stream:%w", err)
	}

	f, err := os.Open(storeDir + "/" + localFilename)
	if err != nil {
		return fmt.Errorf("can not open file:%w", err)
	}

	sendBuf := make([]byte, DefaultPacketSize)

	for {
		l, err := f.Read(sendBuf)
		if err == io.EOF {
			err = writer.Close()
			if err != nil {
				return fmt.Errorf("close stream failed:%w", err)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("read from local file failed:%w", err)
		}

		if err = writer.Write(sendBuf[:l]); err != nil {
			return fmt.Errorf("write file failed:%w", err)
		}
	}

	return nil
}
