package input

import (
	"bytes"
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"fmt"
	"io"
	"strconv"
)

type LineRecordReader struct {
	fsClient  *SDFSSDK.SDFSClient
	inputInfo *idl.InputSplit

	in  *SDFSSDK.StreamReader
	buf *bytes.Buffer

	lineCounter int
	readSize    int
	inputSize   int
}

func NewLineRecordReader(client *SDFSSDK.SDFSClient, inputInfo *idl.InputSplit) (RecordReader, error) {
	formatter := LineRecordReader{
		fsClient:  client,
		inputInfo: inputInfo,
		buf:       &bytes.Buffer{},
	}

	if err := formatter.initRecordReader(); err != nil {
		return nil, err
	}

	return &formatter, nil
}

func (lrr *LineRecordReader) initRecordReader() error {
	inputInfo := lrr.inputInfo
	offset := inputInfo.InputOffset
	if offset > 0 {
		// may start with partial line, in which case we need to drop first line
		// read and check whether previous char is '\n'
		offset -= 1
	}
	// read more content to guarantee the last line will not be incomplete
	// if input size = whole file size, SDFS will only return whole file, it will be safe
	lrr.inputSize = int(inputInfo.InputSize)
	size := lrr.inputSize
	if size == 0 {
		return fmt.Errorf("invalid input split size:0")
	}
	size += MaxLineLength

	in, err := lrr.fsClient.ReadFileAt(inputInfo.InputFileName, int(offset), size)
	if err != nil {
		return fmt.Errorf("get input stream reader failed:%w", err)
	}
	lrr.in = in

	if offset > 0 {
		err = lrr.dropOneKV()
		if err != nil {
			return fmt.Errorf("check first line failed:%w", err)
		}
	}

	return nil
}

func (lrr *LineRecordReader) readFromInputSplit() error {
	content, err := lrr.in.NextPacket()
	lrr.buf.Write(content)
	return err
}

func (lrr *LineRecordReader) dropOneKV() error {
	_, err := lrr.NextKeyValue()
	return err
}

func (lrr *LineRecordReader) NextKeyValue() (*maple_juice.KV, error) {
	if lrr.buf.Len() == 0 {
		err := lrr.readFromInputSplit()
		if err == io.EOF {
			return nil, io.EOF
		}
		if err != nil {
			return nil, fmt.Errorf("read from input split failed:%w", err)
		}
	}

	if lrr.readSize >= lrr.inputSize {
		return nil, io.EOF
	}

	line, err := lrr.buf.ReadBytes('\n')
	if err != nil {
		// ReadBytes returns err != nil if and only if the returned data does not end in delim.
		if lrr.readSize+len(line) < lrr.inputSize {
			err := lrr.readFromInputSplit()
			if err != nil {
				return nil, fmt.Errorf("pad last line: read from input split failed:%w", err)
			}
		}
	}

	lrr.lineCounter += 1
	lrr.readSize += len(line)
	newKV := maple_juice.KV{
		Key:   strconv.Itoa(lrr.lineCounter),
		Value: line,
	}
	return &newKV, nil
}
