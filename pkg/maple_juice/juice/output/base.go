package output

import (
	"ece428_mp4/pkg/maple_juice"
	"encoding/json"
	"fmt"
	"io"
)

type FormatWriter interface {
	Output(kv *maple_juice.KV) error
}

type FormatWriterType int

const (
	JsonMarshalFormatWriterType = iota + 1
)

type JsonMarshalFormatWriter struct {
	writer io.Writer
}

func NewJsonMarshalFormatWriter(writer io.Writer) FormatWriter {
	fw := JsonMarshalFormatWriter{writer: writer}

	return &fw
}

func (fw *JsonMarshalFormatWriter) Output(kv *maple_juice.KV) error {
	marshaled, err := json.Marshal(kv)
	if err != nil {
		return fmt.Errorf("can not marshal cur kv:%w", err)
	}

	_, err = fw.writer.Write(marshaled)
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}

	return nil
}
