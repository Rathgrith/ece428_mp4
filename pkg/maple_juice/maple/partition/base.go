package partition

import "ece428_mp4/pkg/maple_juice"

type Partitioner interface {
	GetPartitionName(*maple_juice.KV) string
}

type PartitionerType int

const (
	PerKeyPartitionerType = iota + 1
)

type NewPartitionerFunc = func() Partitioner

var NewRecordReaderFuncGetter = map[PartitionerType]NewPartitionerFunc{
	PerKeyPartitionerType: NewPerKeyPartitioner,
}

func GetNewRecordReaderFunc(t PartitionerType) NewPartitionerFunc {
	f, exist := NewRecordReaderFuncGetter[t]
	if !exist {
		return nil
	}
	return f
}

type PerKeyPartitioner struct {
}

func NewPerKeyPartitioner() Partitioner {
	return &PerKeyPartitioner{}
}

func (p *PerKeyPartitioner) GetPartitionName(kv *maple_juice.KV) string {
	return kv.Key
}
