// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.12
// source: master_shard/state.proto

package master_shard

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

type ShardSpecIndex int32

const (
	ShardSpecIndex_ShardID       ShardSpecIndex = 0
	ShardSpecIndex_LabelKey      ShardSpecIndex = 1
	ShardSpecIndex_LabelKeyValue ShardSpecIndex = 2
)

// Enum value maps for ShardSpecIndex.
var (
	ShardSpecIndex_name = map[int32]string{
		0: "ShardID",
		1: "LabelKey",
		2: "LabelKeyValue",
	}
	ShardSpecIndex_value = map[string]int32{
		"ShardID":       0,
		"LabelKey":      1,
		"LabelKeyValue": 2,
	}
)

func (x ShardSpecIndex) Enum() *ShardSpecIndex {
	p := new(ShardSpecIndex)
	*p = x
	return p
}

func (x ShardSpecIndex) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ShardSpecIndex) Descriptor() protoreflect.EnumDescriptor {
	return file_master_shard_state_proto_enumTypes[0].Descriptor()
}

func (ShardSpecIndex) Type() protoreflect.EnumType {
	return &file_master_shard_state_proto_enumTypes[0]
}

func (x ShardSpecIndex) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ShardSpecIndex.Descriptor instead.
func (ShardSpecIndex) EnumDescriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{0}
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeHostId string `protobuf:"bytes,1,opt,name=node_host_id,json=nodeHostId,proto3" json:"node_host_id,omitempty"`
	NodeId     uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Addr       string `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetNodeHostId() string {
	if x != nil {
		return x.NodeHostId
	}
	return ""
}

func (x *Node) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *Node) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type NodeSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateVersion uint64            `protobuf:"varint,1,opt,name=state_version,json=stateVersion,proto3" json:"state_version,omitempty"`
	NodeHostId   string            `protobuf:"bytes,2,opt,name=node_host_id,json=nodeHostId,proto3" json:"node_host_id,omitempty"`
	Labels       map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NodeSpec) Reset() {
	*x = NodeSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeSpec) ProtoMessage() {}

func (x *NodeSpec) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeSpec.ProtoReflect.Descriptor instead.
func (*NodeSpec) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{1}
}

func (x *NodeSpec) GetStateVersion() uint64 {
	if x != nil {
		return x.StateVersion
	}
	return 0
}

func (x *NodeSpec) GetNodeHostId() string {
	if x != nil {
		return x.NodeHostId
	}
	return ""
}

func (x *NodeSpec) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

// StateMachine
type ShardSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateVersion uint64            `protobuf:"varint,1,opt,name=state_version,json=stateVersion,proto3" json:"state_version,omitempty"`
	ShardName    string            `protobuf:"bytes,2,opt,name=shard_name,json=shardName,proto3" json:"shard_name,omitempty"`
	ShardId      uint64            `protobuf:"varint,3,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	ProfileName  string            `protobuf:"bytes,4,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
	Replica      int32             `protobuf:"varint,5,opt,name=replica,proto3" json:"replica,omitempty"`
	Initials     map[string]*Node  `protobuf:"bytes,6,rep,name=initials,proto3" json:"initials,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Nodes        map[string]*Node  `protobuf:"bytes,7,rep,name=nodes,proto3" json:"nodes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	LastNodeId   uint64            `protobuf:"varint,8,opt,name=last_node_id,json=lastNodeId,proto3" json:"last_node_id,omitempty"`
	Labels       map[string]string `protobuf:"bytes,9,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ShardSpec) Reset() {
	*x = ShardSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardSpec) ProtoMessage() {}

func (x *ShardSpec) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardSpec.ProtoReflect.Descriptor instead.
func (*ShardSpec) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{2}
}

func (x *ShardSpec) GetStateVersion() uint64 {
	if x != nil {
		return x.StateVersion
	}
	return 0
}

func (x *ShardSpec) GetShardName() string {
	if x != nil {
		return x.ShardName
	}
	return ""
}

func (x *ShardSpec) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *ShardSpec) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

func (x *ShardSpec) GetReplica() int32 {
	if x != nil {
		return x.Replica
	}
	return 0
}

func (x *ShardSpec) GetInitials() map[string]*Node {
	if x != nil {
		return x.Initials
	}
	return nil
}

func (x *ShardSpec) GetNodes() map[string]*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *ShardSpec) GetLastNodeId() uint64 {
	if x != nil {
		return x.LastNodeId
	}
	return 0
}

func (x *ShardSpec) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type IndexData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index string   `protobuf:"bytes,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Keys  []string `protobuf:"bytes,2,rep,name=Keys,proto3" json:"Keys,omitempty"`
}

func (x *IndexData) Reset() {
	*x = IndexData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexData) ProtoMessage() {}

func (x *IndexData) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexData.ProtoReflect.Descriptor instead.
func (*IndexData) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{3}
}

func (x *IndexData) GetIndex() string {
	if x != nil {
		return x.Index
	}
	return ""
}

func (x *IndexData) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type IndexMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Indices []string `protobuf:"bytes,2,rep,name=Indices,proto3" json:"Indices,omitempty"`
}

func (x *IndexMeta) Reset() {
	*x = IndexMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexMeta) ProtoMessage() {}

func (x *IndexMeta) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexMeta.ProtoReflect.Descriptor instead.
func (*IndexMeta) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{4}
}

func (x *IndexMeta) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *IndexMeta) GetIndices() []string {
	if x != nil {
		return x.Indices
	}
	return nil
}

type Indices struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Indices        map[string]*IndexData `protobuf:"bytes,1,rep,name=indices,proto3" json:"indices,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Meta           map[string]*IndexMeta `protobuf:"bytes,2,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	IndexerVersion uint64                `protobuf:"varint,3,opt,name=indexer_version,json=indexerVersion,proto3" json:"indexer_version,omitempty"`
}

func (x *Indices) Reset() {
	*x = Indices{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Indices) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Indices) ProtoMessage() {}

func (x *Indices) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Indices.ProtoReflect.Descriptor instead.
func (*Indices) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{5}
}

func (x *Indices) GetIndices() map[string]*IndexData {
	if x != nil {
		return x.Indices
	}
	return nil
}

func (x *Indices) GetMeta() map[string]*IndexMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *Indices) GetIndexerVersion() uint64 {
	if x != nil {
		return x.IndexerVersion
	}
	return 0
}

type GroupSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateVersion uint64 `protobuf:"varint,1,opt,name=state_version,json=stateVersion,proto3" json:"state_version,omitempty"`
	GroupName    string `protobuf:"bytes,2,opt,name=group_name,json=groupName,proto3" json:"group_name,omitempty"`
	Size         int32  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	ProfileName  string `protobuf:"bytes,4,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
}

func (x *GroupSpec) Reset() {
	*x = GroupSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupSpec) ProtoMessage() {}

func (x *GroupSpec) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupSpec.ProtoReflect.Descriptor instead.
func (*GroupSpec) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{6}
}

func (x *GroupSpec) GetStateVersion() uint64 {
	if x != nil {
		return x.StateVersion
	}
	return 0
}

func (x *GroupSpec) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

func (x *GroupSpec) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *GroupSpec) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

type StateRoot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// shard spec
	ShardsVersion uint64 `protobuf:"varint,1,opt,name=shards_version,json=shardsVersion,proto3" json:"shards_version,omitempty"`
	LastShardId   uint64 `protobuf:"varint,2,opt,name=last_shard_id,json=lastShardId,proto3" json:"last_shard_id,omitempty"`
	// shardname-shardspec
	Shards       map[string]*ShardSpec `protobuf:"bytes,3,rep,name=shards,proto3" json:"shards,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ShardIndices *Indices              `protobuf:"bytes,4,opt,name=shard_indices,json=shardIndices,proto3" json:"shard_indices,omitempty"`
	// nhid-nodespec
	NodesVersion  uint64                `protobuf:"varint,10,opt,name=nodes_version,json=nodesVersion,proto3" json:"nodes_version,omitempty"`
	Nodes         map[string]*NodeSpec  `protobuf:"bytes,11,rep,name=nodes,proto3" json:"nodes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	GroupsVersion uint64                `protobuf:"varint,20,opt,name=groups_version,json=groupsVersion,proto3" json:"groups_version,omitempty"`
	Groups        map[string]*GroupSpec `protobuf:"bytes,21,rep,name=groups,proto3" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *StateRoot) Reset() {
	*x = StateRoot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_shard_state_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateRoot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateRoot) ProtoMessage() {}

func (x *StateRoot) ProtoReflect() protoreflect.Message {
	mi := &file_master_shard_state_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateRoot.ProtoReflect.Descriptor instead.
func (*StateRoot) Descriptor() ([]byte, []int) {
	return file_master_shard_state_proto_rawDescGZIP(), []int{7}
}

func (x *StateRoot) GetShardsVersion() uint64 {
	if x != nil {
		return x.ShardsVersion
	}
	return 0
}

func (x *StateRoot) GetLastShardId() uint64 {
	if x != nil {
		return x.LastShardId
	}
	return 0
}

func (x *StateRoot) GetShards() map[string]*ShardSpec {
	if x != nil {
		return x.Shards
	}
	return nil
}

func (x *StateRoot) GetShardIndices() *Indices {
	if x != nil {
		return x.ShardIndices
	}
	return nil
}

func (x *StateRoot) GetNodesVersion() uint64 {
	if x != nil {
		return x.NodesVersion
	}
	return 0
}

func (x *StateRoot) GetNodes() map[string]*NodeSpec {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *StateRoot) GetGroupsVersion() uint64 {
	if x != nil {
		return x.GroupsVersion
	}
	return 0
}

func (x *StateRoot) GetGroups() map[string]*GroupSpec {
	if x != nil {
		return x.Groups
	}
	return nil
}

var File_master_shard_state_proto protoreflect.FileDescriptor

var file_master_shard_state_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x22, 0x55, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x20, 0x0a, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x22,
	0xc8, 0x01, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x23, 0x0a, 0x0d,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x20, 0x0a, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73,
	0x74, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a,
	0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xdd, 0x04, 0x0a, 0x09, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a,
	0x0a, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x68, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x12, 0x41, 0x0a, 0x08, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x2e,
	0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x38, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x20, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x6f, 0x64,
	0x65, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x09, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x1a, 0x4f, 0x0a, 0x0d, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x1a, 0x4c, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a,
	0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x35, 0x0a, 0x09, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a,
	0x04, 0x4b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x4b, 0x65, 0x79,
	0x73, 0x22, 0x37, 0x0a, 0x09, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x10,
	0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x22, 0xcc, 0x02, 0x0a, 0x07, 0x49,
	0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x12, 0x3c, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x69, 0x63, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x49,
	0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x69, 0x6e, 0x64,
	0x69, 0x63, 0x65, 0x73, 0x12, 0x33, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x1a, 0x53, 0x0a, 0x0c, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x50, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x86, 0x01, 0x0a, 0x09, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x53, 0x70, 0x65, 0x63, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x22, 0x8c, 0x05, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74,
	0x12, 0x25, 0x0a, 0x0e, 0x73, 0x68, 0x61, 0x72, 0x64, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x73, 0x68, 0x61, 0x72, 0x64, 0x73,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x53, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x06, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x61,
	0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x6f, 0x6f, 0x74, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x06, 0x73, 0x68, 0x61, 0x72, 0x64, 0x73, 0x12, 0x3a, 0x0a, 0x0d, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x5f, 0x69, 0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x49,
	0x6e, 0x64, 0x69, 0x63, 0x65, 0x73, 0x52, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x6e, 0x64,
	0x69, 0x63, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6e, 0x6f, 0x64,
	0x65, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x05, 0x6e, 0x6f, 0x64,
	0x65, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f,
	0x74, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x6e, 0x6f,
	0x64, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x14, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x06, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x6f, 0x6f, 0x74, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x1a, 0x52, 0x0a, 0x0b, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x50, 0x0a, 0x0a, 0x4e,
	0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x70,
	0x65, 0x63, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x52, 0x0a,
	0x0b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x2e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x53, 0x70, 0x65, 0x63, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x2a, 0x3e, 0x0a, 0x0e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x68, 0x61, 0x72, 0x64, 0x49, 0x44, 0x10, 0x00,
	0x12, 0x0c, 0x0a, 0x08, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4b, 0x65, 0x79, 0x10, 0x01, 0x12, 0x11,
	0x0a, 0x0d, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x10,
	0x02, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x4c, 0x69, 0x6c, 0x69, 0x74, 0x68, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x2f, 0x6d, 0x6f, 0x78, 0x61,
	0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_master_shard_state_proto_rawDescOnce sync.Once
	file_master_shard_state_proto_rawDescData = file_master_shard_state_proto_rawDesc
)

func file_master_shard_state_proto_rawDescGZIP() []byte {
	file_master_shard_state_proto_rawDescOnce.Do(func() {
		file_master_shard_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_master_shard_state_proto_rawDescData)
	})
	return file_master_shard_state_proto_rawDescData
}

var file_master_shard_state_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_master_shard_state_proto_msgTypes = make([]protoimpl.MessageInfo, 17)
var file_master_shard_state_proto_goTypes = []interface{}{
	(ShardSpecIndex)(0), // 0: master_shard.ShardSpecIndex
	(*Node)(nil),        // 1: master_shard.Node
	(*NodeSpec)(nil),    // 2: master_shard.NodeSpec
	(*ShardSpec)(nil),   // 3: master_shard.ShardSpec
	(*IndexData)(nil),   // 4: master_shard.IndexData
	(*IndexMeta)(nil),   // 5: master_shard.IndexMeta
	(*Indices)(nil),     // 6: master_shard.Indices
	(*GroupSpec)(nil),   // 7: master_shard.GroupSpec
	(*StateRoot)(nil),   // 8: master_shard.StateRoot
	nil,                 // 9: master_shard.NodeSpec.LabelsEntry
	nil,                 // 10: master_shard.ShardSpec.InitialsEntry
	nil,                 // 11: master_shard.ShardSpec.NodesEntry
	nil,                 // 12: master_shard.ShardSpec.LabelsEntry
	nil,                 // 13: master_shard.Indices.IndicesEntry
	nil,                 // 14: master_shard.Indices.MetaEntry
	nil,                 // 15: master_shard.StateRoot.ShardsEntry
	nil,                 // 16: master_shard.StateRoot.NodesEntry
	nil,                 // 17: master_shard.StateRoot.GroupsEntry
}
var file_master_shard_state_proto_depIdxs = []int32{
	9,  // 0: master_shard.NodeSpec.labels:type_name -> master_shard.NodeSpec.LabelsEntry
	10, // 1: master_shard.ShardSpec.initials:type_name -> master_shard.ShardSpec.InitialsEntry
	11, // 2: master_shard.ShardSpec.nodes:type_name -> master_shard.ShardSpec.NodesEntry
	12, // 3: master_shard.ShardSpec.labels:type_name -> master_shard.ShardSpec.LabelsEntry
	13, // 4: master_shard.Indices.indices:type_name -> master_shard.Indices.IndicesEntry
	14, // 5: master_shard.Indices.meta:type_name -> master_shard.Indices.MetaEntry
	15, // 6: master_shard.StateRoot.shards:type_name -> master_shard.StateRoot.ShardsEntry
	6,  // 7: master_shard.StateRoot.shard_indices:type_name -> master_shard.Indices
	16, // 8: master_shard.StateRoot.nodes:type_name -> master_shard.StateRoot.NodesEntry
	17, // 9: master_shard.StateRoot.groups:type_name -> master_shard.StateRoot.GroupsEntry
	1,  // 10: master_shard.ShardSpec.InitialsEntry.value:type_name -> master_shard.Node
	1,  // 11: master_shard.ShardSpec.NodesEntry.value:type_name -> master_shard.Node
	4,  // 12: master_shard.Indices.IndicesEntry.value:type_name -> master_shard.IndexData
	5,  // 13: master_shard.Indices.MetaEntry.value:type_name -> master_shard.IndexMeta
	3,  // 14: master_shard.StateRoot.ShardsEntry.value:type_name -> master_shard.ShardSpec
	2,  // 15: master_shard.StateRoot.NodesEntry.value:type_name -> master_shard.NodeSpec
	7,  // 16: master_shard.StateRoot.GroupsEntry.value:type_name -> master_shard.GroupSpec
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_master_shard_state_proto_init() }
func file_master_shard_state_proto_init() {
	if File_master_shard_state_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_master_shard_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_master_shard_state_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeSpec); i {
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
		file_master_shard_state_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardSpec); i {
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
		file_master_shard_state_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexData); i {
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
		file_master_shard_state_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexMeta); i {
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
		file_master_shard_state_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Indices); i {
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
		file_master_shard_state_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupSpec); i {
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
		file_master_shard_state_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateRoot); i {
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
			RawDescriptor: file_master_shard_state_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   17,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_master_shard_state_proto_goTypes,
		DependencyIndexes: file_master_shard_state_proto_depIdxs,
		EnumInfos:         file_master_shard_state_proto_enumTypes,
		MessageInfos:      file_master_shard_state_proto_msgTypes,
	}.Build()
	File_master_shard_state_proto = out.File
	file_master_shard_state_proto_rawDesc = nil
	file_master_shard_state_proto_goTypes = nil
	file_master_shard_state_proto_depIdxs = nil
}
