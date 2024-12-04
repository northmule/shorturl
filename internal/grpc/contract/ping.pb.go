// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: shorturl/ping.proto

package contract

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CheckStorageConnectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *CheckStorageConnectResponse) Reset() {
	*x = CheckStorageConnectResponse{}
	mi := &file_shorturl_ping_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckStorageConnectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckStorageConnectResponse) ProtoMessage() {}

func (x *CheckStorageConnectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shorturl_ping_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckStorageConnectResponse.ProtoReflect.Descriptor instead.
func (*CheckStorageConnectResponse) Descriptor() ([]byte, []int) {
	return file_shorturl_ping_proto_rawDescGZIP(), []int{0}
}

func (x *CheckStorageConnectResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_shorturl_ping_proto protoreflect.FileDescriptor

var file_shorturl_ping_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2d, 0x0a, 0x1b, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0x72, 0x0a, 0x0b, 0x50, 0x69, 0x6e,
	0x67, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x63, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x25, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0d,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x07, 0x12, 0x05, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x42, 0x0b, 0x5a,
	0x09, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_shorturl_ping_proto_rawDescOnce sync.Once
	file_shorturl_ping_proto_rawDescData = file_shorturl_ping_proto_rawDesc
)

func file_shorturl_ping_proto_rawDescGZIP() []byte {
	file_shorturl_ping_proto_rawDescOnce.Do(func() {
		file_shorturl_ping_proto_rawDescData = protoimpl.X.CompressGZIP(file_shorturl_ping_proto_rawDescData)
	})
	return file_shorturl_ping_proto_rawDescData
}

var file_shorturl_ping_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shorturl_ping_proto_goTypes = []any{
	(*CheckStorageConnectResponse)(nil), // 0: contract.CheckStorageConnectResponse
	(*empty.Empty)(nil),                 // 1: google.protobuf.Empty
}
var file_shorturl_ping_proto_depIdxs = []int32{
	1, // 0: contract.PingHandler.CheckStorageConnect:input_type -> google.protobuf.Empty
	0, // 1: contract.PingHandler.CheckStorageConnect:output_type -> contract.CheckStorageConnectResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shorturl_ping_proto_init() }
func file_shorturl_ping_proto_init() {
	if File_shorturl_ping_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shorturl_ping_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shorturl_ping_proto_goTypes,
		DependencyIndexes: file_shorturl_ping_proto_depIdxs,
		MessageInfos:      file_shorturl_ping_proto_msgTypes,
	}.Build()
	File_shorturl_ping_proto = out.File
	file_shorturl_ping_proto_rawDesc = nil
	file_shorturl_ping_proto_goTypes = nil
	file_shorturl_ping_proto_depIdxs = nil
}
