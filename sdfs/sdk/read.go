package sdk

import (
	"bufio"
	"bytes"
	"context"
	idl2 "ece428_mp4/sdfs/internal/idl"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
)

type StreamReader struct {
	stream idl2.DataNodeService_ReadFileClient
	buf    bytes.Buffer
}

func NewStreamReader(stream idl2.DataNodeService_ReadFileClient) *StreamReader {
	r := StreamReader{stream: stream}

	return &r
}

func (r *StreamReader) NextPacket() ([]byte, error) {
	resp, err := r.stream.Recv()
	if err != nil || resp.GetCode() == idl2.StatusCode_Error {
		return nil, fmt.Errorf("read failed:%w", err)
	}

	if resp.GetCode() == idl2.StatusCode_ReadCompleted {
		return nil, io.EOF
	}

	return resp.GetContent(), err
}

func (c *SDFSClient) ReadFile(filename string) (*StreamReader, error) {
	return c.ReadFileAt(filename, 0, 0)
}

func (c *SDFSClient) ReadFileAt(filename string, offset, size int) (*StreamReader, error) {
	locResp, err := c.getFileLocations(filename)
	if err != nil {
		return nil, fmt.Errorf("query file location failed:%w", err)
	}

	coordinatorHost := locResp.Coordinator.GetHostname()
	coordinator, err := c.rpcManager.GetDataNodePRCClient(coordinatorHost)
	if err != nil {
		return nil, fmt.Errorf("can not dial data node:%w", err)
	}

	readReq := idl2.ReadFileRequest{Filename: filename, Offset: proto.Int32(int32(offset)), Size: proto.Int32(int32(size))}
	stream, err := coordinator.ReadFile(context.Background(), &readReq)
	if err != nil {
		return nil, fmt.Errorf("read file failed:%w", err)
	}

	return NewStreamReader(stream), nil
}

func (c *SDFSClient) ReadFileToLocal(filename string, offset, size int, storeName string, storeDir string) error {
	reader, err := c.ReadFileAt(filename, offset, size)
	if err != nil {
		return fmt.Errorf("get reader failed:%w", err)
	}

	f, err := os.OpenFile(storeDir+"/"+storeName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("can not create file:%w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	for {
		content, err := reader.NextPacket()
		if err == io.EOF {
			err = writer.Flush()
			if err != nil {
				return fmt.Errorf("can not flush buffered content:%w", err)
			}
			return nil
		}

		if err != nil {
			return err
		}

		if _, err = writer.Write(content); err != nil {
			return fmt.Errorf("can not write to local file:%w", err)
		}
	}
}

func (c *SDFSClient) GetFileToLocal(filename, storeName, storeDir string) error {
	return c.ReadFileToLocal(filename, 0, 0, storeName, storeDir)
}

func (c *SDFSClient) GetFileSize(filename string) (int, error) {
	locResp, err := c.getFileLocations(filename)
	if err != nil {
		return 0, err
	}

	targetHost := locResp.GetCoordinator().GetHostname()
	client, err := c.rpcManager.GetDataNodePRCClient(targetHost)
	if err != nil {
		return 0, fmt.Errorf("can not dial data node (%v):%w", targetHost, err)
	}

	resp, err := client.GetFileSize(context.Background(), &idl2.GetFileSizeRequest{Filename: filename})
	if err != nil {
		return 0, fmt.Errorf("get file size from (%v) failed:%w", targetHost, err)
	}

	return int(resp.GetSize()), nil
}
