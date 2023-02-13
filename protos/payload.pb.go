// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: protos/payload.proto

package __

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

type AuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthRequest) Reset() {
	*x = AuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_payload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRequest) ProtoMessage() {}

func (x *AuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_payload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRequest.ProtoReflect.Descriptor instead.
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return file_protos_payload_proto_rawDescGZIP(), []int{0}
}

type AuthResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AuthResult) Reset() {
	*x = AuthResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_payload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResult) ProtoMessage() {}

func (x *AuthResult) ProtoReflect() protoreflect.Message {
	mi := &file_protos_payload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResult.ProtoReflect.Descriptor instead.
func (*AuthResult) Descriptor() ([]byte, []int) {
	return file_protos_payload_proto_rawDescGZIP(), []int{1}
}

type DataChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *DataChunk) Reset() {
	*x = DataChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_payload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataChunk) ProtoMessage() {}

func (x *DataChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_payload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataChunk.ProtoReflect.Descriptor instead.
func (*DataChunk) Descriptor() ([]byte, []int) {
	return file_protos_payload_proto_rawDescGZIP(), []int{2}
}

func (x *DataChunk) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DataChunk) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

type UploadResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecvSize  int32 `protobuf:"varint,1,opt,name=recv_size,json=recvSize,proto3" json:"recv_size,omitempty"`
	IsSuccess bool  `protobuf:"varint,2,opt,name=is_success,json=isSuccess,proto3" json:"is_success,omitempty"`
}

func (x *UploadResult) Reset() {
	*x = UploadResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_payload_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResult) ProtoMessage() {}

func (x *UploadResult) ProtoReflect() protoreflect.Message {
	mi := &file_protos_payload_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResult.ProtoReflect.Descriptor instead.
func (*UploadResult) Descriptor() ([]byte, []int) {
	return file_protos_payload_proto_rawDescGZIP(), []int{3}
}

func (x *UploadResult) GetRecvSize() int32 {
	if x != nil {
		return x.RecvSize
	}
	return 0
}

func (x *UploadResult) GetIsSuccess() bool {
	if x != nil {
		return x.IsSuccess
	}
	return false
}

type DownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AuthKey string `protobuf:"bytes,2,opt,name=auth_key,json=authKey,proto3" json:"auth_key,omitempty"`
}

func (x *DownloadRequest) Reset() {
	*x = DownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_payload_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadRequest) ProtoMessage() {}

func (x *DownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_payload_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadRequest.ProtoReflect.Descriptor instead.
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return file_protos_payload_proto_rawDescGZIP(), []int{4}
}

func (x *DownloadRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DownloadRequest) GetAuthKey() string {
	if x != nil {
		return x.AuthKey
	}
	return ""
}

var File_protos_payload_proto protoreflect.FileDescriptor

var file_protos_payload_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0c, 0x0a, 0x0a, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x31, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x4a, 0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65, 0x63, 0x76, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x3c, 0x0a, 0x0f, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x4b, 0x65, 0x79,
	0x32, 0x95, 0x01, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x2d, 0x0a, 0x0e,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0a,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x0d, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x28, 0x01, 0x12, 0x2d, 0x0a, 0x0b, 0x53,
	0x65, 0x6e, 0x64, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x10, 0x2e, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x30, 0x01, 0x12, 0x2c, 0x0a, 0x0f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0c, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_payload_proto_rawDescOnce sync.Once
	file_protos_payload_proto_rawDescData = file_protos_payload_proto_rawDesc
)

func file_protos_payload_proto_rawDescGZIP() []byte {
	file_protos_payload_proto_rawDescOnce.Do(func() {
		file_protos_payload_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_payload_proto_rawDescData)
	})
	return file_protos_payload_proto_rawDescData
}

var file_protos_payload_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_protos_payload_proto_goTypes = []interface{}{
	(*AuthRequest)(nil),     // 0: AuthRequest
	(*AuthResult)(nil),      // 1: AuthResult
	(*DataChunk)(nil),       // 2: DataChunk
	(*UploadResult)(nil),    // 3: UploadResult
	(*DownloadRequest)(nil), // 4: DownloadRequest
}
var file_protos_payload_proto_depIdxs = []int32{
	2, // 0: Payload.ReceivePayload:input_type -> DataChunk
	4, // 1: Payload.SendPayload:input_type -> DownloadRequest
	0, // 2: Payload.RequestDownload:input_type -> AuthRequest
	3, // 3: Payload.ReceivePayload:output_type -> UploadResult
	2, // 4: Payload.SendPayload:output_type -> DataChunk
	1, // 5: Payload.RequestDownload:output_type -> AuthResult
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_payload_proto_init() }
func file_protos_payload_proto_init() {
	if File_protos_payload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_payload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRequest); i {
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
		file_protos_payload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResult); i {
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
		file_protos_payload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataChunk); i {
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
		file_protos_payload_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadResult); i {
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
		file_protos_payload_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadRequest); i {
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
			RawDescriptor: file_protos_payload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_payload_proto_goTypes,
		DependencyIndexes: file_protos_payload_proto_depIdxs,
		MessageInfos:      file_protos_payload_proto_msgTypes,
	}.Build()
	File_protos_payload_proto = out.File
	file_protos_payload_proto_rawDesc = nil
	file_protos_payload_proto_goTypes = nil
	file_protos_payload_proto_depIdxs = nil
}
