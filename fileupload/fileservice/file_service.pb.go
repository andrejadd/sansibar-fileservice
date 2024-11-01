// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: proto/file_service.proto

package fileservice

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

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"` // Unique identifier for the file
	Offset int64  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`              // Offset for resuming upload
	Data   []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`                   // Actual chunk data
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	mi := &file_proto_file_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_proto_file_service_proto_rawDescGZIP(), []int{0}
}

func (x *Chunk) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *Chunk) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Chunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type UploadStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message    string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	LastOffset int64  `protobuf:"varint,3,opt,name=last_offset,json=lastOffset,proto3" json:"last_offset,omitempty"` // The last successfully uploaded offset
}

func (x *UploadStatus) Reset() {
	*x = UploadStatus{}
	mi := &file_proto_file_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadStatus) ProtoMessage() {}

func (x *UploadStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadStatus.ProtoReflect.Descriptor instead.
func (*UploadStatus) Descriptor() ([]byte, []int) {
	return file_proto_file_service_proto_rawDescGZIP(), []int{1}
}

func (x *UploadStatus) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UploadStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UploadStatus) GetLastOffset() int64 {
	if x != nil {
		return x.LastOffset
	}
	return 0
}

type ResumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId  string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`    // File ID to resume upload
	FileExt string `protobuf:"bytes,2,opt,name=file_ext,json=fileExt,proto3" json:"file_ext,omitempty"` // The file extension
}

func (x *ResumeRequest) Reset() {
	*x = ResumeRequest{}
	mi := &file_proto_file_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResumeRequest) ProtoMessage() {}

func (x *ResumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResumeRequest.ProtoReflect.Descriptor instead.
func (*ResumeRequest) Descriptor() ([]byte, []int) {
	return file_proto_file_service_proto_rawDescGZIP(), []int{2}
}

func (x *ResumeRequest) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *ResumeRequest) GetFileExt() string {
	if x != nil {
		return x.FileExt
	}
	return ""
}

type ResumeStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastOffset int64 `protobuf:"varint,1,opt,name=last_offset,json=lastOffset,proto3" json:"last_offset,omitempty"` // The last uploaded offset
}

func (x *ResumeStatus) Reset() {
	*x = ResumeStatus{}
	mi := &file_proto_file_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResumeStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResumeStatus) ProtoMessage() {}

func (x *ResumeStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResumeStatus.ProtoReflect.Descriptor instead.
func (*ResumeStatus) Descriptor() ([]byte, []int) {
	return file_proto_file_service_proto_rawDescGZIP(), []int{3}
}

func (x *ResumeStatus) GetLastOffset() int64 {
	if x != nil {
		return x.LastOffset
	}
	return 0
}

var File_proto_file_service_proto protoreflect.FileDescriptor

var file_proto_file_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x4c, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x63, 0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x6c, 0x61, 0x73, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x43, 0x0a, 0x0d, 0x52, 0x65,
	0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x65, 0x78, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x22,
	0x2f, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x32, 0x91, 0x01, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x39, 0x0a, 0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x2e, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x19,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x28, 0x01, 0x12, 0x47, 0x0a, 0x0e, 0x49,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1a, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x6a, 0x2f, 0x73, 0x61, 0x6e, 0x73, 0x69, 0x62,
	0x61, 0x72, 0x2d, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_file_service_proto_rawDescOnce sync.Once
	file_proto_file_service_proto_rawDescData = file_proto_file_service_proto_rawDesc
)

func file_proto_file_service_proto_rawDescGZIP() []byte {
	file_proto_file_service_proto_rawDescOnce.Do(func() {
		file_proto_file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_file_service_proto_rawDescData)
	})
	return file_proto_file_service_proto_rawDescData
}

var file_proto_file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_file_service_proto_goTypes = []any{
	(*Chunk)(nil),         // 0: fileservice.Chunk
	(*UploadStatus)(nil),  // 1: fileservice.UploadStatus
	(*ResumeRequest)(nil), // 2: fileservice.ResumeRequest
	(*ResumeStatus)(nil),  // 3: fileservice.ResumeStatus
}
var file_proto_file_service_proto_depIdxs = []int32{
	0, // 0: fileservice.FileService.Upload:input_type -> fileservice.Chunk
	2, // 1: fileservice.FileService.InitiateUpload:input_type -> fileservice.ResumeRequest
	1, // 2: fileservice.FileService.Upload:output_type -> fileservice.UploadStatus
	3, // 3: fileservice.FileService.InitiateUpload:output_type -> fileservice.ResumeStatus
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_file_service_proto_init() }
func file_proto_file_service_proto_init() {
	if File_proto_file_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_file_service_proto_goTypes,
		DependencyIndexes: file_proto_file_service_proto_depIdxs,
		MessageInfos:      file_proto_file_service_proto_msgTypes,
	}.Build()
	File_proto_file_service_proto = out.File
	file_proto_file_service_proto_rawDesc = nil
	file_proto_file_service_proto_goTypes = nil
	file_proto_file_service_proto_depIdxs = nil
}
