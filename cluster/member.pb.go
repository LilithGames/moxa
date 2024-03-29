// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.12
// source: cluster/member.proto

package cluster

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

type MemberType int32

const (
	MemberType_Dragonboat MemberType = 0
)

// Enum value maps for MemberType.
var (
	MemberType_name = map[int32]string{
		0: "Dragonboat",
	}
	MemberType_value = map[string]int32{
		"Dragonboat": 0,
	}
)

func (x MemberType) Enum() *MemberType {
	p := new(MemberType)
	*p = x
	return p
}

func (x MemberType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MemberType) Descriptor() protoreflect.EnumDescriptor {
	return file_cluster_member_proto_enumTypes[0].Descriptor()
}

func (MemberType) Type() protoreflect.EnumType {
	return &file_cluster_member_proto_enumTypes[0]
}

func (x MemberType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MemberType.Descriptor instead.
func (MemberType) EnumDescriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{0}
}

type MemberMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HostName         string     `protobuf:"bytes,1,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	NodeHostId       string     `protobuf:"bytes,2,opt,name=node_host_id,json=nodeHostId,proto3" json:"node_host_id,omitempty"`
	RaftAddress      string     `protobuf:"bytes,3,opt,name=raft_address,json=raftAddress,proto3" json:"raft_address,omitempty"`
	NodeHostIndex    int32      `protobuf:"varint,4,opt,name=node_host_index,json=nodeHostIndex,proto3" json:"node_host_index,omitempty"`
	MasterNodeId     uint64     `protobuf:"varint,5,opt,name=master_node_id,json=masterNodeId,proto3" json:"master_node_id,omitempty"`
	StartupTimestamp int64      `protobuf:"varint,6,opt,name=startup_timestamp,json=startupTimestamp,proto3" json:"startup_timestamp,omitempty"`
	Type             MemberType `protobuf:"varint,7,opt,name=type,proto3,enum=cluster.MemberType" json:"type,omitempty"`
}

func (x *MemberMeta) Reset() {
	*x = MemberMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberMeta) ProtoMessage() {}

func (x *MemberMeta) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberMeta.ProtoReflect.Descriptor instead.
func (*MemberMeta) Descriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{0}
}

func (x *MemberMeta) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *MemberMeta) GetNodeHostId() string {
	if x != nil {
		return x.NodeHostId
	}
	return ""
}

func (x *MemberMeta) GetRaftAddress() string {
	if x != nil {
		return x.RaftAddress
	}
	return ""
}

func (x *MemberMeta) GetNodeHostIndex() int32 {
	if x != nil {
		return x.NodeHostIndex
	}
	return 0
}

func (x *MemberMeta) GetMasterNodeId() uint64 {
	if x != nil {
		return x.MasterNodeId
	}
	return 0
}

func (x *MemberMeta) GetStartupTimestamp() int64 {
	if x != nil {
		return x.StartupTimestamp
	}
	return 0
}

func (x *MemberMeta) GetType() MemberType {
	if x != nil {
		return x.Type
	}
	return MemberType_Dragonboat
}

type MemberShard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShardId    uint64            `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	NodeId     uint64            `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Nodes      map[uint64]string `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	IsLeader   bool              `protobuf:"varint,5,opt,name=is_leader,json=isLeader,proto3" json:"is_leader,omitempty"`
	IsObserver bool              `protobuf:"varint,6,opt,name=is_observer,json=isObserver,proto3" json:"is_observer,omitempty"`
	IsWitness  bool              `protobuf:"varint,7,opt,name=is_witness,json=isWitness,proto3" json:"is_witness,omitempty"`
	Pending    bool              `protobuf:"varint,8,opt,name=pending,proto3" json:"pending,omitempty"`
}

func (x *MemberShard) Reset() {
	*x = MemberShard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_member_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberShard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberShard) ProtoMessage() {}

func (x *MemberShard) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_member_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberShard.ProtoReflect.Descriptor instead.
func (*MemberShard) Descriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{1}
}

func (x *MemberShard) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *MemberShard) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *MemberShard) GetNodes() map[uint64]string {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *MemberShard) GetIsLeader() bool {
	if x != nil {
		return x.IsLeader
	}
	return false
}

func (x *MemberShard) GetIsObserver() bool {
	if x != nil {
		return x.IsObserver
	}
	return false
}

func (x *MemberShard) GetIsWitness() bool {
	if x != nil {
		return x.IsWitness
	}
	return false
}

func (x *MemberShard) GetPending() bool {
	if x != nil {
		return x.Pending
	}
	return false
}

type MemberState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeHostId string                  `protobuf:"bytes,1,opt,name=node_host_id,json=nodeHostId,proto3" json:"node_host_id,omitempty"`
	Meta       *MemberMeta             `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	Shards     map[uint64]*MemberShard `protobuf:"bytes,3,rep,name=shards,proto3" json:"shards,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	LogShards  map[uint64]uint64       `protobuf:"bytes,4,rep,name=log_shards,json=logShards,proto3" json:"log_shards,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *MemberState) Reset() {
	*x = MemberState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_member_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberState) ProtoMessage() {}

func (x *MemberState) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_member_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberState.ProtoReflect.Descriptor instead.
func (*MemberState) Descriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{2}
}

func (x *MemberState) GetNodeHostId() string {
	if x != nil {
		return x.NodeHostId
	}
	return ""
}

func (x *MemberState) GetMeta() *MemberMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *MemberState) GetShards() map[uint64]*MemberShard {
	if x != nil {
		return x.Shards
	}
	return nil
}

func (x *MemberState) GetLogShards() map[uint64]uint64 {
	if x != nil {
		return x.LogShards
	}
	return nil
}

type MemberNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ShardId *uint64 `protobuf:"varint,2,opt,name=shard_id,json=shardId,proto3,oneof" json:"shard_id,omitempty"`
}

func (x *MemberNotify) Reset() {
	*x = MemberNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_member_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberNotify) ProtoMessage() {}

func (x *MemberNotify) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_member_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberNotify.ProtoReflect.Descriptor instead.
func (*MemberNotify) Descriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{3}
}

func (x *MemberNotify) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MemberNotify) GetShardId() uint64 {
	if x != nil && x.ShardId != nil {
		return *x.ShardId
	}
	return 0
}

type MemberGlobalState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Members map[string]*MemberState `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Version uint64                  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *MemberGlobalState) Reset() {
	*x = MemberGlobalState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_member_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberGlobalState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberGlobalState) ProtoMessage() {}

func (x *MemberGlobalState) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_member_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberGlobalState.ProtoReflect.Descriptor instead.
func (*MemberGlobalState) Descriptor() ([]byte, []int) {
	return file_cluster_member_proto_rawDescGZIP(), []int{4}
}

func (x *MemberGlobalState) GetMembers() map[string]*MemberState {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *MemberGlobalState) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_cluster_member_proto protoreflect.FileDescriptor

var file_cluster_member_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x22,
	0x92, 0x02, 0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x1b,
	0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x61, 0x66, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x48,
	0x6f, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0c, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x2b,
	0x0a, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x27, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0xa9, 0x02, 0x0a, 0x0b, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b,
	0x69, 0x73, 0x5f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0a, 0x69, 0x73, 0x4f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x69, 0x73, 0x5f, 0x77, 0x69, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x09, 0x69, 0x73, 0x57, 0x69, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70,
	0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x1a, 0x38, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0xe5, 0x02, 0x0a, 0x0b, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x20, 0x0a, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x27, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x38, 0x0a, 0x06, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x73, 0x12, 0x42, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x4c, 0x6f, 0x67, 0x53, 0x68, 0x61, 0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09,
	0x6c, 0x6f, 0x67, 0x53, 0x68, 0x61, 0x72, 0x64, 0x73, 0x1a, 0x4f, 0x0a, 0x0b, 0x53, 0x68, 0x61,
	0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x68, 0x61, 0x72, 0x64, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3c, 0x0a, 0x0e, 0x4c, 0x6f,
	0x67, 0x53, 0x68, 0x61, 0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4f, 0x0a, 0x0c, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x08,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00,
	0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09,
	0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x22, 0xc2, 0x01, 0x0a, 0x11, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x41, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x50, 0x0a, 0x0c,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x1c,
	0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a,
	0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x74, 0x10, 0x00, 0x42, 0x25, 0x5a, 0x23,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4c, 0x69, 0x6c, 0x69, 0x74,
	0x68, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x2f, 0x6d, 0x6f, 0x78, 0x61, 0x2f, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_member_proto_rawDescOnce sync.Once
	file_cluster_member_proto_rawDescData = file_cluster_member_proto_rawDesc
)

func file_cluster_member_proto_rawDescGZIP() []byte {
	file_cluster_member_proto_rawDescOnce.Do(func() {
		file_cluster_member_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_member_proto_rawDescData)
	})
	return file_cluster_member_proto_rawDescData
}

var file_cluster_member_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cluster_member_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_cluster_member_proto_goTypes = []interface{}{
	(MemberType)(0),           // 0: cluster.MemberType
	(*MemberMeta)(nil),        // 1: cluster.MemberMeta
	(*MemberShard)(nil),       // 2: cluster.MemberShard
	(*MemberState)(nil),       // 3: cluster.MemberState
	(*MemberNotify)(nil),      // 4: cluster.MemberNotify
	(*MemberGlobalState)(nil), // 5: cluster.MemberGlobalState
	nil,                       // 6: cluster.MemberShard.NodesEntry
	nil,                       // 7: cluster.MemberState.ShardsEntry
	nil,                       // 8: cluster.MemberState.LogShardsEntry
	nil,                       // 9: cluster.MemberGlobalState.MembersEntry
}
var file_cluster_member_proto_depIdxs = []int32{
	0, // 0: cluster.MemberMeta.type:type_name -> cluster.MemberType
	6, // 1: cluster.MemberShard.nodes:type_name -> cluster.MemberShard.NodesEntry
	1, // 2: cluster.MemberState.meta:type_name -> cluster.MemberMeta
	7, // 3: cluster.MemberState.shards:type_name -> cluster.MemberState.ShardsEntry
	8, // 4: cluster.MemberState.log_shards:type_name -> cluster.MemberState.LogShardsEntry
	9, // 5: cluster.MemberGlobalState.members:type_name -> cluster.MemberGlobalState.MembersEntry
	2, // 6: cluster.MemberState.ShardsEntry.value:type_name -> cluster.MemberShard
	3, // 7: cluster.MemberGlobalState.MembersEntry.value:type_name -> cluster.MemberState
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_cluster_member_proto_init() }
func file_cluster_member_proto_init() {
	if File_cluster_member_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_member_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberMeta); i {
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
		file_cluster_member_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberShard); i {
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
		file_cluster_member_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberState); i {
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
		file_cluster_member_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberNotify); i {
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
		file_cluster_member_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberGlobalState); i {
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
	file_cluster_member_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cluster_member_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cluster_member_proto_goTypes,
		DependencyIndexes: file_cluster_member_proto_depIdxs,
		EnumInfos:         file_cluster_member_proto_enumTypes,
		MessageInfos:      file_cluster_member_proto_msgTypes,
	}.Build()
	File_cluster_member_proto = out.File
	file_cluster_member_proto_rawDesc = nil
	file_cluster_member_proto_goTypes = nil
	file_cluster_member_proto_depIdxs = nil
}
