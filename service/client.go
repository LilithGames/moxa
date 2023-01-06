package service

import (
	"fmt"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/connectivity"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"github.com/LilithGames/moxa/cluster"
)

type IClient interface {
	Close() error
	Wait(ctx context.Context, timeout time.Duration) error
	NodeHost() NodeHostClient
	Spec() SpecClient
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(cm cluster.Manager, target string) (IClient, error) {
	builder := NewServiceResolverBuilder(cm)

	retry_opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponentialWithJitter(100 * time.Millisecond, 0.01)),
		grpc_retry.WithCodes(RouteNotFound),
	}

	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"DragonboatRoundRobinBalancer"}`),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retry_opts...)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial err: %w", err)
	}
	return &Client{conn}, nil
}

func (it *Client) NodeHost() NodeHostClient {
	return NewNodeHostClientWrapper(NewNodeHostClient(it.conn))
}

func (it *Client) Spec() SpecClient {
	return NewSpecClientWrapper(NewSpecClient(it.conn))
}

func (it *Client) Wait(ctx context.Context, timeout time.Duration) error {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		it.conn.Connect()
		state := it.conn.GetState()
		log.Println("[INFO]", fmt.Sprintf("conn: %v", state))
		if state == connectivity.Ready {
			return nil
		} else if state == connectivity.Shutdown {
			return fmt.Errorf("connectivity.Shutdown")
		}
		if !it.conn.WaitForStateChange(cctx, state) {
			return fmt.Errorf("WaitForStateChange %w", cctx.Err())
		}
	}
}

func (it *Client) Close() error {
	return it.conn.Close()
}


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
	return it.NodeHostClient.CreateSnapshot(WithRouteShardLeader(ctx, in.ShardId), in , opts...)
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
