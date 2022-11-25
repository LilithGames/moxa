package service

import (
	"fmt"
	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"

	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/cluster"
)

type SpecService struct {
	cm cluster.Manager
	UnimplementedSpecServer
}

func RegisterSpecService(cs *ClusterService) {
	svc := &SpecService{cm: cs.ClusterManager()}
	RegisterSpecServer(cs.GrpcServiceRegistrar(), svc)
	cs.AddGrpcGatewayRegister(RegisterSpecHandlerFromEndpoint)
}

func (it *SpecService) AddShardSpec(ctx context.Context, req *AddShardSpecRequest) (*AddShardSpecResponse, error) {
	nodes := make([]*master_shard.NodeCreateView, 0, it.cm.Members().Nums())
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		nodes = append(nodes, &master_shard.NodeCreateView{NodeHostId: m.Meta.NodeHostId, Addr: m.Meta.RaftAddress})
		return true
	})
	resp, err := it.cm.Client().MasterShard().CreateShard(ctx, &master_shard.CreateShardRequest{Name: req.ShardName, Nodes: nodes})
	if err != nil {
		return &AddShardSpecResponse{}, wrapDragonboatErr(fmt.Errorf("Client().MasterShard().CreateShard(%s) err: %w", req.ShardName, err))
	}
	return &AddShardSpecResponse{Shard: ShardSpecSiri(resp.Shard)}, nil
}

func (it *SpecService) RemoveShardSpec(ctx context.Context, req *RemoveShardSpecRequest) (*RemoveShardSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().DeleteShard(ctx, &master_shard.DeleteShardRequest{Name: req.ShardName})
	if err != nil {
		return &RemoveShardSpecResponse{}, wrapDragonboatErr(fmt.Errorf("Client().MasterShard().DeleteShard(%s) err: %w", req.ShardName, err))
	}
	return &RemoveShardSpecResponse{Updated: resp.Updated}, nil
}

func (it *SpecService) RebalanceShardSpec(ctx context.Context, req *RebalanceShardSpecRequest) (*RebalanceShardSpecResponse, error) {
	nodes := make([]*master_shard.NodeCreateView, 0, it.cm.Members().Nums())
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		nodes = append(nodes, &master_shard.NodeCreateView{NodeHostId: m.Meta.NodeHostId, Addr: m.Meta.RaftAddress})
		return true
	})
	resp, err := it.cm.Client().MasterShard().UpdateShard(ctx, &master_shard.UpdateShardRequest{Name: req.ShardName, Nodes: nodes})
	if err != nil {
		return &RebalanceShardSpecResponse{}, wrapDragonboatErr(fmt.Errorf("Client().MasterShard().UpdateShard(%s) err: %w", req.ShardName, err))
	}
	return &RebalanceShardSpecResponse{Updated: resp.Updated}, nil
}

func (it *SpecService) ListShardSpec(ctx context.Context, req *ListShardSpecRequest) (*ListShardSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().ListShards(ctx, &master_shard.ListShardsRequest{})
	if err != nil {
		return &ListShardSpecResponse{}, wrapDragonboatErr(fmt.Errorf("Client().MasterShard().ListShards() err: %w", err))
	}
	return &ListShardSpecResponse{Shards: lo.Map(resp.Shards, WithNopIndex(ShardSpecSiri)), Version: resp.StateVersion}, nil
}

func wrapDragonboatErr(err error) error {
	if err == nil {
		return nil
	}
	if 400 <= runtime.GetDragonboatErrorCode(err) && runtime.GetDragonboatErrorCode(err) < 500 {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}
