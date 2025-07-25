// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: antibruteforce.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CheckRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Login         string                 `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Ip            string                 `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	mi := &file_antibruteforce_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{0}
}

func (x *CheckRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *CheckRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *CheckRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type CheckResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ok            bool                   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	mi := &file_antibruteforce_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type ResetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Login         string                 `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Ip            string                 `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResetRequest) Reset() {
	*x = ResetRequest{}
	mi := &file_antibruteforce_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetRequest) ProtoMessage() {}

func (x *ResetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetRequest.ProtoReflect.Descriptor instead.
func (*ResetRequest) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{2}
}

func (x *ResetRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *ResetRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type NetworkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Network       string                 `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NetworkRequest) Reset() {
	*x = NetworkRequest{}
	mi := &file_antibruteforce_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NetworkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkRequest) ProtoMessage() {}

func (x *NetworkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkRequest.ProtoReflect.Descriptor instead.
func (*NetworkRequest) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{3}
}

func (x *NetworkRequest) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

var File_antibruteforce_proto protoreflect.FileDescriptor

const file_antibruteforce_proto_rawDesc = "" +
	"\n" +
	"\x14antibruteforce.proto\x12\x0eantibruteforce\x1a\x1bgoogle/protobuf/empty.proto\"P\n" +
	"\fCheckRequest\x12\x14\n" +
	"\x05login\x18\x01 \x01(\tR\x05login\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x0e\n" +
	"\x02ip\x18\x03 \x01(\tR\x02ip\"\x1f\n" +
	"\rCheckResponse\x12\x0e\n" +
	"\x02ok\x18\x01 \x01(\bR\x02ok\"4\n" +
	"\fResetRequest\x12\x14\n" +
	"\x05login\x18\x01 \x01(\tR\x05login\x12\x0e\n" +
	"\x02ip\x18\x02 \x01(\tR\x02ip\"*\n" +
	"\x0eNetworkRequest\x12\x18\n" +
	"\anetwork\x18\x01 \x01(\tR\anetwork2\xc7\x03\n" +
	"\x0eAntiBruteforce\x12D\n" +
	"\x05Check\x12\x1c.antibruteforce.CheckRequest\x1a\x1d.antibruteforce.CheckResponse\x12=\n" +
	"\x05Reset\x12\x1c.antibruteforce.ResetRequest\x1a\x16.google.protobuf.Empty\x12H\n" +
	"\x0eAddToBlacklist\x12\x1e.antibruteforce.NetworkRequest\x1a\x16.google.protobuf.Empty\x12M\n" +
	"\x13RemoveFromBlacklist\x12\x1e.antibruteforce.NetworkRequest\x1a\x16.google.protobuf.Empty\x12H\n" +
	"\x0eAddToWhitelist\x12\x1e.antibruteforce.NetworkRequest\x1a\x16.google.protobuf.Empty\x12M\n" +
	"\x13RemoveFromWhitelist\x12\x1e.antibruteforce.NetworkRequest\x1a\x16.google.protobuf.EmptyB\aZ\x05./;pbb\x06proto3"

var (
	file_antibruteforce_proto_rawDescOnce sync.Once
	file_antibruteforce_proto_rawDescData []byte
)

func file_antibruteforce_proto_rawDescGZIP() []byte {
	file_antibruteforce_proto_rawDescOnce.Do(func() {
		file_antibruteforce_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_antibruteforce_proto_rawDesc), len(file_antibruteforce_proto_rawDesc)))
	})
	return file_antibruteforce_proto_rawDescData
}

var file_antibruteforce_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_antibruteforce_proto_goTypes = []any{
	(*CheckRequest)(nil),   // 0: antibruteforce.CheckRequest
	(*CheckResponse)(nil),  // 1: antibruteforce.CheckResponse
	(*ResetRequest)(nil),   // 2: antibruteforce.ResetRequest
	(*NetworkRequest)(nil), // 3: antibruteforce.NetworkRequest
	(*emptypb.Empty)(nil),  // 4: google.protobuf.Empty
}
var file_antibruteforce_proto_depIdxs = []int32{
	0, // 0: antibruteforce.AntiBruteforce.Check:input_type -> antibruteforce.CheckRequest
	2, // 1: antibruteforce.AntiBruteforce.Reset:input_type -> antibruteforce.ResetRequest
	3, // 2: antibruteforce.AntiBruteforce.AddToBlacklist:input_type -> antibruteforce.NetworkRequest
	3, // 3: antibruteforce.AntiBruteforce.RemoveFromBlacklist:input_type -> antibruteforce.NetworkRequest
	3, // 4: antibruteforce.AntiBruteforce.AddToWhitelist:input_type -> antibruteforce.NetworkRequest
	3, // 5: antibruteforce.AntiBruteforce.RemoveFromWhitelist:input_type -> antibruteforce.NetworkRequest
	1, // 6: antibruteforce.AntiBruteforce.Check:output_type -> antibruteforce.CheckResponse
	4, // 7: antibruteforce.AntiBruteforce.Reset:output_type -> google.protobuf.Empty
	4, // 8: antibruteforce.AntiBruteforce.AddToBlacklist:output_type -> google.protobuf.Empty
	4, // 9: antibruteforce.AntiBruteforce.RemoveFromBlacklist:output_type -> google.protobuf.Empty
	4, // 10: antibruteforce.AntiBruteforce.AddToWhitelist:output_type -> google.protobuf.Empty
	4, // 11: antibruteforce.AntiBruteforce.RemoveFromWhitelist:output_type -> google.protobuf.Empty
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_antibruteforce_proto_init() }
func file_antibruteforce_proto_init() {
	if File_antibruteforce_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_antibruteforce_proto_rawDesc), len(file_antibruteforce_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_antibruteforce_proto_goTypes,
		DependencyIndexes: file_antibruteforce_proto_depIdxs,
		MessageInfos:      file_antibruteforce_proto_msgTypes,
	}.Build()
	File_antibruteforce_proto = out.File
	file_antibruteforce_proto_goTypes = nil
	file_antibruteforce_proto_depIdxs = nil
}
