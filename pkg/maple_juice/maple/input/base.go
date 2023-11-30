package input

import (
	"ece428_mp4/idl"
	"ece428_mp4/pkg/maple_juice"
	SDFSSDK "ece428_mp4/sdfs/sdk"
)

type RecordReaderType int

const (
	LineRecordReaderType = iota + 1
)

type RecordReader interface {
	NextKeyValue() (*maple_juice.KV, error)
}

type NewRecordReaderFunc = func(*SDFSSDK.SDFSClient, *idl.InputSplit) (RecordReader, error)

const (
	MaxLineLength = 1024 * 32
)

var NewRecordReaderFuncGetter = map[RecordReaderType]NewRecordReaderFunc{
	LineRecordReaderType: NewLineRecordReader,
}

func GetNewRecordReaderFunc(t RecordReaderType) NewRecordReaderFunc {
	f, exist := NewRecordReaderFuncGetter[t]
	if !exist {
		return nil
	}
	return f
}
