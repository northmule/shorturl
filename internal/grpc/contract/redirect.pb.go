// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: shorturl/redirect.proto

package contract

import (
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

type RedirectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RedirectRequest) Reset() {
	*x = RedirectRequest{}
	mi := &file_shorturl_redirect_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedirectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectRequest) ProtoMessage() {}

func (x *RedirectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shorturl_redirect_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectRequest.ProtoReflect.Descriptor instead.
func (*RedirectRequest) Descriptor() ([]byte, []int) {
	return file_shorturl_redirect_proto_rawDescGZIP(), []int{0}
}

func (x *RedirectRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RedirectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *RedirectResponse) Reset() {
	*x = RedirectResponse{}
	mi := &file_shorturl_redirect_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedirectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectResponse) ProtoMessage() {}

func (x *RedirectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shorturl_redirect_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectResponse.ProtoReflect.Descriptor instead.
func (*RedirectResponse) Descriptor() ([]byte, []int) {
	return file_shorturl_redirect_proto_rawDescGZIP(), []int{1}
}

func (x *RedirectResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_shorturl_redirect_proto protoreflect.FileDescriptor

var file_shorturl_redirect_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x21, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x24, 0x0a, 0x10, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x32, 0x61, 0x0a, 0x0f, 0x52, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x4e, 0x0a,
	0x08, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e,
	0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x0b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x05, 0x12, 0x03, 0x2f, 0x69, 0x64, 0x42, 0x0b, 0x5a,
	0x09, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_shorturl_redirect_proto_rawDescOnce sync.Once
	file_shorturl_redirect_proto_rawDescData = file_shorturl_redirect_proto_rawDesc
)

func file_shorturl_redirect_proto_rawDescGZIP() []byte {
	file_shorturl_redirect_proto_rawDescOnce.Do(func() {
		file_shorturl_redirect_proto_rawDescData = protoimpl.X.CompressGZIP(file_shorturl_redirect_proto_rawDescData)
	})
	return file_shorturl_redirect_proto_rawDescData
}

var file_shorturl_redirect_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_shorturl_redirect_proto_goTypes = []any{
	(*RedirectRequest)(nil),  // 0: contract.RedirectRequest
	(*RedirectResponse)(nil), // 1: contract.RedirectResponse
}
var file_shorturl_redirect_proto_depIdxs = []int32{
	0, // 0: contract.RedirectHandler.Redirect:input_type -> contract.RedirectRequest
	1, // 1: contract.RedirectHandler.Redirect:output_type -> contract.RedirectResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shorturl_redirect_proto_init() }
func file_shorturl_redirect_proto_init() {
	if File_shorturl_redirect_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shorturl_redirect_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shorturl_redirect_proto_goTypes,
		DependencyIndexes: file_shorturl_redirect_proto_depIdxs,
		MessageInfos:      file_shorturl_redirect_proto_msgTypes,
	}.Build()
	File_shorturl_redirect_proto = out.File
	file_shorturl_redirect_proto_rawDesc = nil
	file_shorturl_redirect_proto_goTypes = nil
	file_shorturl_redirect_proto_depIdxs = nil
}
