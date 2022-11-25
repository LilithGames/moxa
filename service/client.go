package service

import (
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/LilithGames/moxa/cluster"
)

type IClient interface {
	Close() error
	NodeHost() NodeHostClient
	Spec() SpecClient
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(cm cluster.Manager) (IClient, error) {
	builder := NewServiceResolverBuilder(cm)
	conn, err := grpc.Dial("dragonboat://:8001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"DragonboatRoundRobinBalancer"}`),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial err: %w", err)
	}
	return &Client{conn}, nil
}

func (it *Client) NodeHost() NodeHostClient {
	return NewShardClientWrapper(NewNodeHostClient(it.conn))
}

func (it *Client) Spec() SpecClient {
	return NewSpecClient(it.conn)
}

func (it *Client) Close() error {
	return it.conn.Close()
}


type NodeHostClientWrapper struct {
	NodeHostClient
}

func NewShardClientWrapper(c NodeHostClient) NodeHostClient {
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
func (it *NodeHostClientWrapper) ListNode(ctx context.Context, in *ShardListNodeRequest, opts ...grpc.CallOption) (*ShardListNodeResponse, error) {
	return it.NodeHostClient.ListNode(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) TransferLeader(ctx context.Context, in *ShardTransferLeaderRequest, opts ...grpc.CallOption) (*ShardTransferLeaderResponse, error) {
	return it.NodeHostClient.TransferLeader(WithRouteShard(ctx, in.ShardId), in, opts...)
}
func (it *NodeHostClientWrapper) CreateSnapshot(ctx context.Context, in *CreateSnapshotRequest, opts ...grpc.CallOption) (*CreateSnapshotResponse, error) {
	return it.NodeHostClient.CreateSnapshot(WithRouteShardLeader(ctx, in.ShardId), in , opts...)
}
