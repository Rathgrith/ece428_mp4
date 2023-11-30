// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: maple_juice.proto

package idl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InputFormatterType int32

const (
	InputFormatterType_Undefined     InputFormatterType = 0
	InputFormatterType_LineTextInput InputFormatterType = 1
)

// Enum value maps for InputFormatterType.
var (
	InputFormatterType_name = map[int32]string{
		0: "Undefined",
		1: "LineTextInput",
	}
	InputFormatterType_value = map[string]int32{
		"Undefined":     0,
		"LineTextInput": 1,
	}
)

func (x InputFormatterType) Enum() *InputFormatterType {
	p := new(InputFormatterType)
	*p = x
	return p
}

func (x InputFormatterType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InputFormatterType) Descriptor() protoreflect.EnumDescriptor {
	return file_maple_juice_proto_enumTypes[0].Descriptor()
}

func (InputFormatterType) Type() protoreflect.EnumType {
	return &file_maple_juice_proto_enumTypes[0]
}

func (x InputFormatterType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InputFormatterType.Descriptor instead.
func (InputFormatterType) EnumDescriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{0}
}

type InputSplit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InputFileName string `protobuf:"bytes,1,opt,name=input_file_name,json=inputFileName,proto3" json:"input_file_name,omitempty"`
	InputOffset   int32  `protobuf:"varint,2,opt,name=input_offset,json=inputOffset,proto3" json:"input_offset,omitempty"`
	InputSize     int32  `protobuf:"varint,3,opt,name=input_size,json=inputSize,proto3" json:"input_size,omitempty"`
}

func (x *InputSplit) Reset() {
	*x = InputSplit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_maple_juice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputSplit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputSplit) ProtoMessage() {}

func (x *InputSplit) ProtoReflect() protoreflect.Message {
	mi := &file_maple_juice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputSplit.ProtoReflect.Descriptor instead.
func (*InputSplit) Descriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{0}
}

func (x *InputSplit) GetInputFileName() string {
	if x != nil {
		return x.InputFileName
	}
	return ""
}

func (x *InputSplit) GetInputOffset() int32 {
	if x != nil {
		return x.InputOffset
	}
	return 0
}

func (x *InputSplit) GetInputSize() int32 {
	if x != nil {
		return x.InputSize
	}
	return 0
}

type RunMapleTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExeName                    string        `protobuf:"bytes,1,opt,name=exe_name,json=exeName,proto3" json:"exe_name,omitempty"`
	IntermediateFilenamePrefix string        `protobuf:"bytes,2,opt,name=intermediate_filename_prefix,json=intermediateFilenamePrefix,proto3" json:"intermediate_filename_prefix,omitempty"`
	InputSplits                []*InputSplit `protobuf:"bytes,3,rep,name=input_splits,json=inputSplits,proto3" json:"input_splits,omitempty"`
	AttemptId                  string        `protobuf:"bytes,4,opt,name=attempt_id,json=attemptId,proto3" json:"attempt_id,omitempty"`
}

func (x *RunMapleTaskRequest) Reset() {
	*x = RunMapleTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_maple_juice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunMapleTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunMapleTaskRequest) ProtoMessage() {}

func (x *RunMapleTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_maple_juice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunMapleTaskRequest.ProtoReflect.Descriptor instead.
func (*RunMapleTaskRequest) Descriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{1}
}

func (x *RunMapleTaskRequest) GetExeName() string {
	if x != nil {
		return x.ExeName
	}
	return ""
}

func (x *RunMapleTaskRequest) GetIntermediateFilenamePrefix() string {
	if x != nil {
		return x.IntermediateFilenamePrefix
	}
	return ""
}

func (x *RunMapleTaskRequest) GetInputSplits() []*InputSplit {
	if x != nil {
		return x.InputSplits
	}
	return nil
}

func (x *RunMapleTaskRequest) GetAttemptId() string {
	if x != nil {
		return x.AttemptId
	}
	return ""
}

type RunMapleTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code                 StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=idl.StatusCode" json:"code,omitempty"`
	TmpIntermediateFiles []string   `protobuf:"bytes,2,rep,name=tmp_intermediate_files,json=tmpIntermediateFiles,proto3" json:"tmp_intermediate_files,omitempty"`
	ErrInfo              *string    `protobuf:"bytes,3,opt,name=ErrInfo,proto3,oneof" json:"ErrInfo,omitempty"`
}

func (x *RunMapleTaskResponse) Reset() {
	*x = RunMapleTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_maple_juice_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunMapleTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunMapleTaskResponse) ProtoMessage() {}

func (x *RunMapleTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_maple_juice_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunMapleTaskResponse.ProtoReflect.Descriptor instead.
func (*RunMapleTaskResponse) Descriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{2}
}

func (x *RunMapleTaskResponse) GetCode() StatusCode {
	if x != nil {
		return x.Code
	}
	return StatusCode_Unknown
}

func (x *RunMapleTaskResponse) GetTmpIntermediateFiles() []string {
	if x != nil {
		return x.TmpIntermediateFiles
	}
	return nil
}

func (x *RunMapleTaskResponse) GetErrInfo() string {
	if x != nil && x.ErrInfo != nil {
		return *x.ErrInfo
	}
	return ""
}

type ExecuteMapleJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExeName                    string   `protobuf:"bytes,1,opt,name=exe_name,json=exeName,proto3" json:"exe_name,omitempty"`
	IntermediateFilenamePrefix string   `protobuf:"bytes,2,opt,name=intermediate_filename_prefix,json=intermediateFilenamePrefix,proto3" json:"intermediate_filename_prefix,omitempty"`
	InputFiles                 []string `protobuf:"bytes,3,rep,name=input_files,json=inputFiles,proto3" json:"input_files,omitempty"`
	NumMaples                  int32    `protobuf:"varint,4,opt,name=num_maples,json=numMaples,proto3" json:"num_maples,omitempty"`
}

func (x *ExecuteMapleJobRequest) Reset() {
	*x = ExecuteMapleJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_maple_juice_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteMapleJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteMapleJobRequest) ProtoMessage() {}

func (x *ExecuteMapleJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_maple_juice_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteMapleJobRequest.ProtoReflect.Descriptor instead.
func (*ExecuteMapleJobRequest) Descriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{3}
}

func (x *ExecuteMapleJobRequest) GetExeName() string {
	if x != nil {
		return x.ExeName
	}
	return ""
}

func (x *ExecuteMapleJobRequest) GetIntermediateFilenamePrefix() string {
	if x != nil {
		return x.IntermediateFilenamePrefix
	}
	return ""
}

func (x *ExecuteMapleJobRequest) GetInputFiles() []string {
	if x != nil {
		return x.InputFiles
	}
	return nil
}

func (x *ExecuteMapleJobRequest) GetNumMaples() int32 {
	if x != nil {
		return x.NumMaples
	}
	return 0
}

type ExecuteMapleJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=idl.StatusCode" json:"code,omitempty"`
}

func (x *ExecuteMapleJobResponse) Reset() {
	*x = ExecuteMapleJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_maple_juice_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteMapleJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteMapleJobResponse) ProtoMessage() {}

func (x *ExecuteMapleJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_maple_juice_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteMapleJobResponse.ProtoReflect.Descriptor instead.
func (*ExecuteMapleJobResponse) Descriptor() ([]byte, []int) {
	return file_maple_juice_proto_rawDescGZIP(), []int{4}
}

func (x *ExecuteMapleJobResponse) GetCode() StatusCode {
	if x != nil {
		return x.Code
	}
	return StatusCode_Unknown
}

var File_maple_juice_proto protoreflect.FileDescriptor

var file_maple_juice_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x61, 0x70, 0x6c, 0x65, 0x5f, 0x6a, 0x75, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x69, 0x64, 0x6c, 0x1a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x70,
	0x6c, 0x69, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x22, 0xc5, 0x01,
	0x0a, 0x13, 0x52, 0x75, 0x6e, 0x4d, 0x61, 0x70, 0x6c, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x78, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x78, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x40, 0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x65,
	0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66,
	0x69, 0x78, 0x12, 0x32, 0x0a, 0x0c, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x73, 0x70, 0x6c, 0x69,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x52, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x53, 0x70, 0x6c, 0x69, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x49, 0x64, 0x22, 0x9c, 0x01, 0x0a, 0x14, 0x52, 0x75, 0x6e, 0x4d, 0x61, 0x70,
	0x6c, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x69,
	0x64, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x16, 0x74, 0x6d, 0x70, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x14, 0x74, 0x6d, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x07, 0x45, 0x72, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x45, 0x72,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x45, 0x72, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x22, 0xb5, 0x01, 0x0a, 0x16, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65,
	0x4d, 0x61, 0x70, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x65, 0x78, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x65, 0x78, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x1c, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x1f, 0x0a, 0x0b,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x0a,
	0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x4d, 0x61, 0x70, 0x6c, 0x65, 0x73, 0x22, 0x3e, 0x0a, 0x17,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x4d, 0x61, 0x70, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x2a, 0x36, 0x0a, 0x12,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x74, 0x65, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10,
	0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x69, 0x6e, 0x65, 0x54, 0x65, 0x78, 0x74, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x10, 0x01, 0x32, 0x58, 0x0a, 0x11, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0c, 0x52, 0x75, 0x6e,
	0x4d, 0x61, 0x70, 0x6c, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x18, 0x2e, 0x69, 0x64, 0x6c, 0x2e,
	0x52, 0x75, 0x6e, 0x4d, 0x61, 0x70, 0x6c, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x52, 0x75, 0x6e, 0x4d, 0x61, 0x70,
	0x6c, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x07,
	0x5a, 0x05, 0x2e, 0x2f, 0x69, 0x64, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_maple_juice_proto_rawDescOnce sync.Once
	file_maple_juice_proto_rawDescData = file_maple_juice_proto_rawDesc
)

func file_maple_juice_proto_rawDescGZIP() []byte {
	file_maple_juice_proto_rawDescOnce.Do(func() {
		file_maple_juice_proto_rawDescData = protoimpl.X.CompressGZIP(file_maple_juice_proto_rawDescData)
	})
	return file_maple_juice_proto_rawDescData
}

var file_maple_juice_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_maple_juice_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_maple_juice_proto_goTypes = []interface{}{
	(InputFormatterType)(0),         // 0: idl.InputFormatterType
	(*InputSplit)(nil),              // 1: idl.InputSplit
	(*RunMapleTaskRequest)(nil),     // 2: idl.RunMapleTaskRequest
	(*RunMapleTaskResponse)(nil),    // 3: idl.RunMapleTaskResponse
	(*ExecuteMapleJobRequest)(nil),  // 4: idl.ExecuteMapleJobRequest
	(*ExecuteMapleJobResponse)(nil), // 5: idl.ExecuteMapleJobResponse
	(StatusCode)(0),                 // 6: idl.StatusCode
}
var file_maple_juice_proto_depIdxs = []int32{
	1, // 0: idl.RunMapleTaskRequest.input_splits:type_name -> idl.InputSplit
	6, // 1: idl.RunMapleTaskResponse.code:type_name -> idl.StatusCode
	6, // 2: idl.ExecuteMapleJobResponse.code:type_name -> idl.StatusCode
	2, // 3: idl.NodeManageService.RunMapleTask:input_type -> idl.RunMapleTaskRequest
	3, // 4: idl.NodeManageService.RunMapleTask:output_type -> idl.RunMapleTaskResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_maple_juice_proto_init() }
func file_maple_juice_proto_init() {
	if File_maple_juice_proto != nil {
		return
	}
	file_share_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_maple_juice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputSplit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_maple_juice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunMapleTaskRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_maple_juice_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunMapleTaskResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_maple_juice_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecuteMapleJobRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_maple_juice_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecuteMapleJobResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_maple_juice_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_maple_juice_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_maple_juice_proto_goTypes,
		DependencyIndexes: file_maple_juice_proto_depIdxs,
		EnumInfos:         file_maple_juice_proto_enumTypes,
		MessageInfos:      file_maple_juice_proto_msgTypes,
	}.Build()
	File_maple_juice_proto = out.File
	file_maple_juice_proto_rawDesc = nil
	file_maple_juice_proto_goTypes = nil
	file_maple_juice_proto_depIdxs = nil
}