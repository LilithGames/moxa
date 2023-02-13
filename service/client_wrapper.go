package service

import (
	"context"

	"google.golang.org/grpc"
)

type NodeHostClientWrapper struct {
	NodeHostClient
}

func NewNodeHostClientWrapper(c NodeHostClient) NodeHostClient {
	return &NodeHostClientWrapper{NodeHostClient: c}
}
func (it *NodeHostClientWrapper) ShardHealthz(ctx context.Context, in *ShardHealthzRequest, opts ...grpc.CallOption) (*ShardHealthzResponse, error) {
	return it.NodeHostClient.ShardHealthz(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) ShardNop(ctx context.Context, in *ShardNopRequest, opts ...grpc.CallOption) (*ShardNopResponse, error) {
	return it.NodeHostClient.ShardNop(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) AddNode(ctx context.Context, in *ShardAddNodeRequest, opts ...grpc.CallOption) (*ShardAddNodeResponse, error) {
	return it.NodeHostClient.AddNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) RemoveNode(ctx context.Context, in *ShardRemoveNodeRequest, opts ...grpc.CallOption) (*ShardRemoveNodeResponse, error) {
	return it.NodeHostClient.RemoveNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) SyncAddNode(ctx context.Context, in *SyncAddNodeRequest, opts ...grpc.CallOption) (*SyncAddNodeResponse, error) {
	return it.NodeHostClient.SyncAddNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) SyncRemoveNode(ctx context.Context, in *SyncRemoveNodeRequest, opts ...grpc.CallOption) (*SyncRemoveNodeResponse, error) {
	return it.NodeHostClient.SyncRemoveNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) ListNode(ctx context.Context, in *ShardListNodeRequest, opts ...grpc.CallOption) (*ShardListNodeResponse, error) {
	return it.NodeHostClient.ListNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) TransferLeader(ctx context.Context, in *ShardTransferLeaderRequest, opts ...grpc.CallOption) (*ShardTransferLeaderResponse, error) {
	return it.NodeHostClient.TransferLeader(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) CreateSnapshot(ctx context.Context, in *CreateSnapshotRequest, opts ...grpc.CallOption) (*CreateSnapshotResponse, error) {
	return it.NodeHostClient.CreateSnapshot(WithRouteShardLeader(ctx, in.ShardId), in, opts...)
}

type SpecClientWrapper struct {
	SpecClient
}

func NewSpecClientWrapper(c SpecClient) SpecClient {
	return &SpecClientWrapper{c}
}

func (it *SpecClientWrapper) CreateMigration(ctx context.Context, in *CreateMigrationRequest, opts ...grpc.CallOption) (*CreateMigrationResponse, error) {
	return it.SpecClient.CreateMigration(WithRouteShard(ctx, in.ShardId), in, opts...)
}

func (it *SpecClientWrapper) GetMigration(ctx context.Context, in *QueryMigrationRequest, opts ...grpc.CallOption) (*QueryMigrationResponse, error) {
	return it.SpecClient.GetMigration(WithRouteShard(ctx, in.ShardId), in, opts...)
}
