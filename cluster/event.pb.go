// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.2
// source: cluster/event.proto

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

type EventTopic int32

const (
	EventTopic_None                      EventTopic = 0
	EventTopic_NodeHostNodeReady         EventTopic = 1000
	EventTopic_NodeHostMembershipChanged EventTopic = 1001
	EventTopic_NodeHostLeaderUpdated     EventTopic = 1002
	EventTopic_ShardSpecUpdating         EventTopic = 2000
	EventTopic_ShardSpecChanging         EventTopic = 2001
)

// Enum value maps for EventTopic.
var (
	EventTopic_name = map[int32]string{
		0:    "None",
		1000: "NodeHostNodeReady",
		1001: "NodeHostMembershipChanged",
		1002: "NodeHostLeaderUpdated",
		2000: "ShardSpecUpdating",
		2001: "ShardSpecChanging",
	}
	EventTopic_value = map[string]int32{
		"None":                      0,
		"NodeHostNodeReady":         1000,
		"NodeHostMembershipChanged": 1001,
		"NodeHostLeaderUpdated":     1002,
		"ShardSpecUpdating":         2000,
		"ShardSpecChanging":         2001,
	}
)

func (x EventTopic) Enum() *EventTopic {
	p := new(EventTopic)
	*p = x
	return p
}

func (x EventTopic) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventTopic) Descriptor() protoreflect.EnumDescriptor {
	return file_cluster_event_proto_enumTypes[0].Descriptor()
}

func (EventTopic) Type() protoreflect.EnumType {
	return &file_cluster_event_proto_enumTypes[0]
}

func (x EventTopic) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventTopic.Descriptor instead.
func (EventTopic) EnumDescriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{0}
}

type ShardSpecChangingType int32

const (
	ShardSpecChangingType_Adding      ShardSpecChangingType = 0
	ShardSpecChangingType_Deleting    ShardSpecChangingType = 1
	ShardSpecChangingType_NodeJoining ShardSpecChangingType = 2
	ShardSpecChangingType_NodeLeaving ShardSpecChangingType = 3
	ShardSpecChangingType_NodeMoving  ShardSpecChangingType = 4
	ShardSpecChangingType_Cleanup     ShardSpecChangingType = 5
)

// Enum value maps for ShardSpecChangingType.
var (
	ShardSpecChangingType_name = map[int32]string{
		0: "Adding",
		1: "Deleting",
		2: "NodeJoining",
		3: "NodeLeaving",
		4: "NodeMoving",
		5: "Cleanup",
	}
	ShardSpecChangingType_value = map[string]int32{
		"Adding":      0,
		"Deleting":    1,
		"NodeJoining": 2,
		"NodeLeaving": 3,
		"NodeMoving":  4,
		"Cleanup":     5,
	}
)

func (x ShardSpecChangingType) Enum() *ShardSpecChangingType {
	p := new(ShardSpecChangingType)
	*p = x
	return p
}

func (x ShardSpecChangingType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ShardSpecChangingType) Descriptor() protoreflect.EnumDescriptor {
	return file_cluster_event_proto_enumTypes[1].Descriptor()
}

func (ShardSpecChangingType) Type() protoreflect.EnumType {
	return &file_cluster_event_proto_enumTypes[1]
}

func (x ShardSpecChangingType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ShardSpecChangingType.Descriptor instead.
func (ShardSpecChangingType) EnumDescriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{1}
}

type NodeHostNodeReadyEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShardId uint64 `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	NodeId  uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *NodeHostNodeReadyEvent) Reset() {
	*x = NodeHostNodeReadyEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeHostNodeReadyEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeHostNodeReadyEvent) ProtoMessage() {}

func (x *NodeHostNodeReadyEvent) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeHostNodeReadyEvent.ProtoReflect.Descriptor instead.
func (*NodeHostNodeReadyEvent) Descriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{0}
}

func (x *NodeHostNodeReadyEvent) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *NodeHostNodeReadyEvent) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

type NodeHostMembershipChangedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShardId uint64 `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	NodeId  uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *NodeHostMembershipChangedEvent) Reset() {
	*x = NodeHostMembershipChangedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeHostMembershipChangedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeHostMembershipChangedEvent) ProtoMessage() {}

func (x *NodeHostMembershipChangedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeHostMembershipChangedEvent.ProtoReflect.Descriptor instead.
func (*NodeHostMembershipChangedEvent) Descriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{1}
}

func (x *NodeHostMembershipChangedEvent) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *NodeHostMembershipChangedEvent) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

type NodeHostLeaderUpdatedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShardId  uint64 `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	NodeId   uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Term     uint64 `protobuf:"varint,3,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId uint64 `protobuf:"varint,4,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
}

func (x *NodeHostLeaderUpdatedEvent) Reset() {
	*x = NodeHostLeaderUpdatedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeHostLeaderUpdatedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeHostLeaderUpdatedEvent) ProtoMessage() {}

func (x *NodeHostLeaderUpdatedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeHostLeaderUpdatedEvent.ProtoReflect.Descriptor instead.
func (*NodeHostLeaderUpdatedEvent) Descriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{2}
}

func (x *NodeHostLeaderUpdatedEvent) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *NodeHostLeaderUpdatedEvent) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *NodeHostLeaderUpdatedEvent) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *NodeHostLeaderUpdatedEvent) GetLeaderId() uint64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

type ShardSpecUpdatingEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateVersion uint64 `protobuf:"varint,1,opt,name=state_version,json=stateVersion,proto3" json:"state_version,omitempty"`
}

func (x *ShardSpecUpdatingEvent) Reset() {
	*x = ShardSpecUpdatingEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardSpecUpdatingEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardSpecUpdatingEvent) ProtoMessage() {}

func (x *ShardSpecUpdatingEvent) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardSpecUpdatingEvent.ProtoReflect.Descriptor instead.
func (*ShardSpecUpdatingEvent) Descriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{3}
}

func (x *ShardSpecUpdatingEvent) GetStateVersion() uint64 {
	if x != nil {
		return x.StateVersion
	}
	return 0
}

type ShardSpecChangingEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type           ShardSpecChangingType `protobuf:"varint,1,opt,name=type,proto3,enum=cluster.ShardSpecChangingType" json:"type,omitempty"`
	ShardId        uint64                `protobuf:"varint,2,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	PreviousNodes  map[uint64]string     `protobuf:"bytes,3,rep,name=previous_nodes,json=previousNodes,proto3" json:"previous_nodes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CurrentNodes   map[uint64]string     `protobuf:"bytes,4,rep,name=current_nodes,json=currentNodes,proto3" json:"current_nodes,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	PreviousNodeId *uint64               `protobuf:"varint,5,opt,name=previous_node_id,json=previousNodeId,proto3,oneof" json:"previous_node_id,omitempty"`
	CurrentNodeId  *uint64               `protobuf:"varint,6,opt,name=current_node_id,json=currentNodeId,proto3,oneof" json:"current_node_id,omitempty"`
}

func (x *ShardSpecChangingEvent) Reset() {
	*x = ShardSpecChangingEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cluster_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardSpecChangingEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardSpecChangingEvent) ProtoMessage() {}

func (x *ShardSpecChangingEvent) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardSpecChangingEvent.ProtoReflect.Descriptor instead.
func (*ShardSpecChangingEvent) Descriptor() ([]byte, []int) {
	return file_cluster_event_proto_rawDescGZIP(), []int{4}
}

func (x *ShardSpecChangingEvent) GetType() ShardSpecChangingType {
	if x != nil {
		return x.Type
	}
	return ShardSpecChangingType_Adding
}

func (x *ShardSpecChangingEvent) GetShardId() uint64 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *ShardSpecChangingEvent) GetPreviousNodes() map[uint64]string {
	if x != nil {
		return x.PreviousNodes
	}
	return nil
}

func (x *ShardSpecChangingEvent) GetCurrentNodes() map[uint64]string {
	if x != nil {
		return x.CurrentNodes
	}
	return nil
}

func (x *ShardSpecChangingEvent) GetPreviousNodeId() uint64 {
	if x != nil && x.PreviousNodeId != nil {
		return *x.PreviousNodeId
	}
	return 0
}

func (x *ShardSpecChangingEvent) GetCurrentNodeId() uint64 {
	if x != nil && x.CurrentNodeId != nil {
		return *x.CurrentNodeId
	}
	return 0
}

var File_cluster_event_proto protoreflect.FileDescriptor

var file_cluster_event_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x22, 0x4c,
	0x0a, 0x16, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x61, 0x64, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x54, 0x0a, 0x1e,
	0x4e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68,
	0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x19,
	0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x22, 0x81, 0x01, 0x0a, 0x1a, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e,
	0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3d, 0x0a, 0x16, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53,
	0x70, 0x65, 0x63, 0x55, 0x70, 0x64, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xa2, 0x04, 0x0a, 0x16, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53,
	0x70, 0x65, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x32, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e,
	0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70,
	0x65, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12,
	0x59, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x69, 0x6e, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75,
	0x73, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x56, 0x0a, 0x0d, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x31, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x68, 0x61, 0x72,
	0x64, 0x53, 0x70, 0x65, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64,
	0x65, 0x73, 0x12, 0x2d, 0x0a, 0x10, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0e,
	0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x2b, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x48, 0x01, 0x52, 0x0d, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x88, 0x01, 0x01, 0x1a, 0x40,
	0x0a, 0x12, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x3f, 0x0a, 0x11, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x2a, 0x9a, 0x01, 0x0a, 0x0a, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e,
	0x65, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x11, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4e,
	0x6f, 0x64, 0x65, 0x52, 0x65, 0x61, 0x64, 0x79, 0x10, 0xe8, 0x07, 0x12, 0x1e, 0x0a, 0x19, 0x4e,
	0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69,
	0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x10, 0xe9, 0x07, 0x12, 0x1a, 0x0a, 0x15, 0x4e,
	0x6f, 0x64, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x10, 0xea, 0x07, 0x12, 0x16, 0x0a, 0x11, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x53, 0x70, 0x65, 0x63, 0x55, 0x70, 0x64, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x10, 0xd0, 0x0f, 0x12,
	0x16, 0x0a, 0x11, 0x53, 0x68, 0x61, 0x72, 0x64, 0x53, 0x70, 0x65, 0x63, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x69, 0x6e, 0x67, 0x10, 0xd1, 0x0f, 0x2a, 0x70, 0x0a, 0x15, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x53, 0x70, 0x65, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0a, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x6f,
	0x64, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4e,
	0x6f, 0x64, 0x65, 0x4c, 0x65, 0x61, 0x76, 0x69, 0x6e, 0x67, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a,
	0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x6c, 0x65, 0x61, 0x6e, 0x75, 0x70, 0x10, 0x05, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4c, 0x69, 0x6c, 0x69, 0x74, 0x68, 0x47, 0x61,
	0x6d, 0x65, 0x73, 0x2f, 0x6d, 0x6f, 0x78, 0x61, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_event_proto_rawDescOnce sync.Once
	file_cluster_event_proto_rawDescData = file_cluster_event_proto_rawDesc
)

func file_cluster_event_proto_rawDescGZIP() []byte {
	file_cluster_event_proto_rawDescOnce.Do(func() {
		file_cluster_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_event_proto_rawDescData)
	})
	return file_cluster_event_proto_rawDescData
}

var file_cluster_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_cluster_event_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_cluster_event_proto_goTypes = []interface{}{
	(EventTopic)(0),                        // 0: cluster.EventTopic
	(ShardSpecChangingType)(0),             // 1: cluster.ShardSpecChangingType
	(*NodeHostNodeReadyEvent)(nil),         // 2: cluster.NodeHostNodeReadyEvent
	(*NodeHostMembershipChangedEvent)(nil), // 3: cluster.NodeHostMembershipChangedEvent
	(*NodeHostLeaderUpdatedEvent)(nil),     // 4: cluster.NodeHostLeaderUpdatedEvent
	(*ShardSpecUpdatingEvent)(nil),         // 5: cluster.ShardSpecUpdatingEvent
	(*ShardSpecChangingEvent)(nil),         // 6: cluster.ShardSpecChangingEvent
	nil,                                    // 7: cluster.ShardSpecChangingEvent.PreviousNodesEntry
	nil,                                    // 8: cluster.ShardSpecChangingEvent.CurrentNodesEntry
}
var file_cluster_event_proto_depIdxs = []int32{
	1, // 0: cluster.ShardSpecChangingEvent.type:type_name -> cluster.ShardSpecChangingType
	7, // 1: cluster.ShardSpecChangingEvent.previous_nodes:type_name -> cluster.ShardSpecChangingEvent.PreviousNodesEntry
	8, // 2: cluster.ShardSpecChangingEvent.current_nodes:type_name -> cluster.ShardSpecChangingEvent.CurrentNodesEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_cluster_event_proto_init() }
func file_cluster_event_proto_init() {
	if File_cluster_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cluster_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeHostNodeReadyEvent); i {
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
		file_cluster_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeHostMembershipChangedEvent); i {
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
		file_cluster_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeHostLeaderUpdatedEvent); i {
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
		file_cluster_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardSpecUpdatingEvent); i {
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
		file_cluster_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardSpecChangingEvent); i {
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
	file_cluster_event_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cluster_event_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cluster_event_proto_goTypes,
		DependencyIndexes: file_cluster_event_proto_depIdxs,
		EnumInfos:         file_cluster_event_proto_enumTypes,
		MessageInfos:      file_cluster_event_proto_msgTypes,
	}.Build()
	File_cluster_event_proto = out.File
	file_cluster_event_proto_rawDesc = nil
	file_cluster_event_proto_goTypes = nil
	file_cluster_event_proto_depIdxs = nil
}