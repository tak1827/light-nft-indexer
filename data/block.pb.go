// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: data/block.proto

package data

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlockType int32

const (
	BlockType_BLOCK_TYPE_UNSPECIFIED               BlockType = 0
	BlockType_BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED  BlockType = 1
	BlockType_BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED BlockType = 2
)

// Enum value maps for BlockType.
var (
	BlockType_name = map[int32]string{
		0: "BLOCK_TYPE_UNSPECIFIED",
		1: "BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED",
		2: "BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED",
	}
	BlockType_value = map[string]int32{
		"BLOCK_TYPE_UNSPECIFIED":               0,
		"BLOCK_TYPE_LAST_FACTORY_LOG_FETCHED":  1,
		"BLOCK_TYPE_LAST_TRANSFER_LOG_FETCHED": 2,
	}
)

func (x BlockType) Enum() *BlockType {
	p := new(BlockType)
	*p = x
	return p
}

func (x BlockType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BlockType) Descriptor() protoreflect.EnumDescriptor {
	return file_data_block_proto_enumTypes[0].Descriptor()
}

func (BlockType) Type() protoreflect.EnumType {
	return &file_data_block_proto_enumTypes[0]
}

func (x BlockType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BlockType.Descriptor instead.
func (BlockType) EnumDescriptor() ([]byte, []int) {
	return file_data_block_proto_rawDescGZIP(), []int{0}
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height        uint64                 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Hash          string                 `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Time          uint64                 `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Type          BlockType              `protobuf:"varint,6,opt,name=type,proto3,enum=tak1827.lightnftindexer.data.BlockType" json:"type,omitempty"`
	SubIdentifier string                 `protobuf:"bytes,7,opt,name=sub_identifier,json=subIdentifier,proto3" json:"sub_identifier,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_data_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_data_block_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Block) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Block) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Block) GetType() BlockType {
	if x != nil {
		return x.Type
	}
	return BlockType_BLOCK_TYPE_UNSPECIFIED
}

func (x *Block) GetSubIdentifier() string {
	if x != nil {
		return x.SubIdentifier
	}
	return ""
}

func (x *Block) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Block) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_data_block_proto protoreflect.FileDescriptor

var file_data_block_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1c, 0x74, 0x61, 0x6b, 0x31, 0x38, 0x32, 0x37, 0x2e, 0x6c, 0x69, 0x67, 0x68,
	0x74, 0x6e, 0x66, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xad, 0x02, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x74, 0x61, 0x6b, 0x31,
	0x38, 0x32, 0x37, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x66, 0x74, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x65, 0x72, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x75, 0x62, 0x5f,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x73, 0x75, 0x62, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x05, 0x4a, 0x04, 0x08, 0x05, 0x10,
	0x06, 0x2a, 0x7a, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a,
	0x0a, 0x16, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x27, 0x0a, 0x23, 0x42, 0x4c,
	0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x41, 0x53, 0x54, 0x5f, 0x46, 0x41,
	0x43, 0x54, 0x4f, 0x52, 0x59, 0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x46, 0x45, 0x54, 0x43, 0x48, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x28, 0x0a, 0x24, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4c, 0x41, 0x53, 0x54, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x45, 0x52, 0x5f,
	0x4c, 0x4f, 0x47, 0x5f, 0x46, 0x45, 0x54, 0x43, 0x48, 0x45, 0x44, 0x10, 0x02, 0x42, 0x2d, 0x48,
	0x02, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x61,
	0x6b, 0x31, 0x38, 0x32, 0x37, 0x2f, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x2d, 0x6e, 0x66, 0x74, 0x2d,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_block_proto_rawDescOnce sync.Once
	file_data_block_proto_rawDescData = file_data_block_proto_rawDesc
)

func file_data_block_proto_rawDescGZIP() []byte {
	file_data_block_proto_rawDescOnce.Do(func() {
		file_data_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_block_proto_rawDescData)
	})
	return file_data_block_proto_rawDescData
}

var file_data_block_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_data_block_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_data_block_proto_goTypes = []interface{}{
	(BlockType)(0),                // 0: tak1827.lightnftindexer.data.BlockType
	(*Block)(nil),                 // 1: tak1827.lightnftindexer.data.Block
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_data_block_proto_depIdxs = []int32{
	0, // 0: tak1827.lightnftindexer.data.Block.type:type_name -> tak1827.lightnftindexer.data.BlockType
	2, // 1: tak1827.lightnftindexer.data.Block.created_at:type_name -> google.protobuf.Timestamp
	2, // 2: tak1827.lightnftindexer.data.Block.updated_at:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_data_block_proto_init() }
func file_data_block_proto_init() {
	if File_data_block_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
			RawDescriptor: file_data_block_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_block_proto_goTypes,
		DependencyIndexes: file_data_block_proto_depIdxs,
		EnumInfos:         file_data_block_proto_enumTypes,
		MessageInfos:      file_data_block_proto_msgTypes,
	}.Build()
	File_data_block_proto = out.File
	file_data_block_proto_rawDesc = nil
	file_data_block_proto_goTypes = nil
	file_data_block_proto_depIdxs = nil
}
