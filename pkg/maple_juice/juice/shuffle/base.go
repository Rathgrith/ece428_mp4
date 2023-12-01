package shuffle

import (
	"fmt"
	"hash/crc32"
)

type IntermediateShuffler interface {
	Shuffle() ([][]string, error)
}

type ShufflerType int

const (
	HashShufflerType = iota + 1
)

type ShufflerCreateFunc = func([]string, int) IntermediateShuffler

var ShufflerCreateFuncGetter = map[ShufflerType]ShufflerCreateFunc{
	HashShufflerType: NewHashShuffler,
}

func GetShufflerCreateFunc(t ShufflerType) ShufflerCreateFunc {
	f, exist := ShufflerCreateFuncGetter[t]
	if !exist {
		return nil
	}
	return f
}

type HashShuffler struct {
	intermediateFilenames []string
	num                   int
}

func NewHashShuffler(files []string, num int) IntermediateShuffler {
	shuffler := HashShuffler{
		intermediateFilenames: files,
		num:                   num,
	}

	return &shuffler
}

func (s *HashShuffler) Shuffle() ([][]string, error) {
	res := make([][]string, s.num)
	if len(s.intermediateFilenames) == 0 {
		return nil, fmt.Errorf("can not shuffle a empty input files")
	}

	for _, filename := range s.intermediateFilenames {
		idx := crc32.ChecksumIEEE([]byte(filename)) % uint32(s.num)
		res[idx] = append(res[idx], filename)
	}

	return res, nil
}
