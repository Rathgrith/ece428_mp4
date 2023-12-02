// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: TaskScheduler.proto

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

type TaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskType    string `protobuf:"bytes,1,opt,name=taskType,proto3" json:"taskType,omitempty"`
	Exe         string `protobuf:"bytes,2,opt,name=exe,proto3" json:"exe,omitempty"`
	NumJobs     int32  `protobuf:"varint,3,opt,name=numJobs,proto3" json:"numJobs,omitempty"`
	Prefix      string `protobuf:"bytes,4,opt,name=prefix,proto3" json:"prefix,omitempty"`
	SrcDir1     string `protobuf:"bytes,5,opt,name=srcDir1,proto3" json:"srcDir1,omitempty"`
	SrcDir2     string `protobuf:"bytes,6,opt,name=srcDir2,proto3" json:"srcDir2,omitempty"`
	Regex       string `protobuf:"bytes,7,opt,name=regex,proto3" json:"regex,omitempty"`
	JoinColumn1 int32  `protobuf:"varint,8,opt,name=joinColumn1,proto3" json:"joinColumn1,omitempty"`
	JoinColumn2 int32  `protobuf:"varint,9,opt,name=joinColumn2,proto3" json:"joinColumn2,omitempty"`
	DestFile    string `protobuf:"bytes,10,opt,name=destFile,proto3" json:"destFile,omitempty"`
}

func (x *TaskRequest) Reset() {
	*x = TaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TaskScheduler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskRequest) ProtoMessage() {}

func (x *TaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_TaskScheduler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskRequest.ProtoReflect.Descriptor instead.
func (*TaskRequest) Descriptor() ([]byte, []int) {
	return file_TaskScheduler_proto_rawDescGZIP(), []int{0}
}

func (x *TaskRequest) GetTaskType() string {
	if x != nil {
		return x.TaskType
	}
	return ""
}

func (x *TaskRequest) GetExe() string {
	if x != nil {
		return x.Exe
	}
	return ""
}

func (x *TaskRequest) GetNumJobs() int32 {
	if x != nil {
		return x.NumJobs
	}
	return 0
}

func (x *TaskRequest) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *TaskRequest) GetSrcDir1() string {
	if x != nil {
		return x.SrcDir1
	}
	return ""
}

func (x *TaskRequest) GetSrcDir2() string {
	if x != nil {
		return x.SrcDir2
	}
	return ""
}

func (x *TaskRequest) GetRegex() string {
	if x != nil {
		return x.Regex
	}
	return ""
}

func (x *TaskRequest) GetJoinColumn1() int32 {
	if x != nil {
		return x.JoinColumn1
	}
	return 0
}

func (x *TaskRequest) GetJoinColumn2() int32 {
	if x != nil {
		return x.JoinColumn2
	}
	return 0
}

func (x *TaskRequest) GetDestFile() string {
	if x != nil {
		return x.DestFile
	}
	return ""
}

type TaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TaskResponse) Reset() {
	*x = TaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TaskScheduler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponse) ProtoMessage() {}

func (x *TaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_TaskScheduler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponse.ProtoReflect.Descriptor instead.
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return file_TaskScheduler_proto_rawDescGZIP(), []int{1}
}

func (x *TaskResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_TaskScheduler_proto protoreflect.FileDescriptor

var file_TaskScheduler_proto_rawDesc = []byte{
	0x0a, 0x13, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x72, 0x22, 0x97, 0x02, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65,
	0x78, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x75, 0x6d, 0x4a, 0x6f, 0x62, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x4a, 0x6f, 0x62, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x72, 0x63, 0x44, 0x69, 0x72, 0x31, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x72, 0x63, 0x44, 0x69, 0x72, 0x31, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x72, 0x63, 0x44, 0x69, 0x72, 0x32, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x72, 0x63, 0x44, 0x69, 0x72, 0x32, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x67, 0x65,
	0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78, 0x12, 0x20,
	0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x31, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x31,
	0x12, 0x20, 0x0a, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x32, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6a, 0x6f, 0x69, 0x6e, 0x43, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x32, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x28,
	0x0a, 0x0c, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x5d, 0x0a, 0x13, 0x4d, 0x61, 0x70, 0x6c,
	0x65, 0x4a, 0x75, 0x69, 0x63, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x12,
	0x46, 0x0a, 0x0b, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x1a,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x54, 0x61, 0x73,
	0x6b, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x69, 0x64, 0x6c,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_TaskScheduler_proto_rawDescOnce sync.Once
	file_TaskScheduler_proto_rawDescData = file_TaskScheduler_proto_rawDesc
)

func file_TaskScheduler_proto_rawDescGZIP() []byte {
	file_TaskScheduler_proto_rawDescOnce.Do(func() {
		file_TaskScheduler_proto_rawDescData = protoimpl.X.CompressGZIP(file_TaskScheduler_proto_rawDescData)
	})
	return file_TaskScheduler_proto_rawDescData
}

var file_TaskScheduler_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_TaskScheduler_proto_goTypes = []interface{}{
	(*TaskRequest)(nil),  // 0: TaskScheduler.TaskRequest
	(*TaskResponse)(nil), // 1: TaskScheduler.TaskResponse
}
var file_TaskScheduler_proto_depIdxs = []int32{
	0, // 0: TaskScheduler.MapleJuiceScheduler.EnqueueTask:input_type -> TaskScheduler.TaskRequest
	1, // 1: TaskScheduler.MapleJuiceScheduler.EnqueueTask:output_type -> TaskScheduler.TaskResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_TaskScheduler_proto_init() }
func file_TaskScheduler_proto_init() {
	if File_TaskScheduler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_TaskScheduler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskRequest); i {
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
		file_TaskScheduler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_TaskScheduler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_TaskScheduler_proto_goTypes,
		DependencyIndexes: file_TaskScheduler_proto_depIdxs,
		MessageInfos:      file_TaskScheduler_proto_msgTypes,
	}.Build()
	File_TaskScheduler_proto = out.File
	file_TaskScheduler_proto_rawDesc = nil
	file_TaskScheduler_proto_goTypes = nil
	file_TaskScheduler_proto_depIdxs = nil
}