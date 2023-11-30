package split

import (
	"ece428_mp4/idl"
	SDFSSDK "ece428_mp4/sdfs/sdk"
	"fmt"
	"math"
)

type InputSplitter interface {
}

type Info struct {
	Split     *idl.InputSplit
	Locations []string
}

type RawFileSplitter struct {
	inputFiles []string
	fsClient   *SDFSSDK.SDFSClient
	nMaples    int

	fileSizes     map[string]int
	totalSize     int64
	fileLocations map[string][]string
	splitsInfo    []*Info
}

func NewRawFileSplitter(inputFiles []string, client *SDFSSDK.SDFSClient, numMaples int) *RawFileSplitter {
	splitter := RawFileSplitter{
		inputFiles: inputFiles,
		fsClient:   client,
		nMaples:    numMaples,

		fileSizes:     make(map[string]int),
		fileLocations: make(map[string][]string),
		splitsInfo:    make([]*Info, 0),
	}

	return &splitter
}

func (s *RawFileSplitter) Split() ([]*Info, error) {
	for _, handleFUnc := range []func() error{
		s.getFileSizes, s.getFileLocations, s.getSplits,
	} {
		if err := handleFUnc(); err != nil {
			return nil, err
		}
	}

	return s.splitsInfo, nil
}

func (s *RawFileSplitter) getFileSizes() error {
	for _, filename := range s.inputFiles {
		size, err := s.fsClient.GetFileSize(filename)
		if err != nil {
			return fmt.Errorf("can not get size of file (%s):%w", filename, err)
		}
		s.fileSizes[filename] = size
		s.totalSize += int64(size)
	}

	return nil
}

func (s *RawFileSplitter) getFileLocations() error {
	for _, filename := range s.inputFiles {
		locs, err := s.fsClient.GetFileLocationHosts(filename)
		if err != nil {
			return fmt.Errorf("can not get Locations of file (%s):%w", filename, err)
		}
		s.fileLocations[filename] = locs
	}

	return nil
}

func (s *RawFileSplitter) getSplits() error {
	// currently simple consider one input file TODO
	filename := s.inputFiles[0]
	fileSize := s.fileSizes[filename]
	expectSplitSize := int(math.Ceil(float64(fileSize / s.nMaples)))
	offset := 0
	for {
		size := expectSplitSize
		if offset+size > fileSize {
			size = fileSize - offset
		}
		split := idl.InputSplit{
			InputFileName: filename,
			InputOffset:   int32(offset),
			InputSize:     int32(expectSplitSize),
		}
		offset += size
		s.splitsInfo = append(s.splitsInfo, &Info{
			Split:     &split,
			Locations: s.fileLocations[filename],
		})
		if offset >= fileSize {
			break
		}
	}

	return nil
}
