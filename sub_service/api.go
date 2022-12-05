package sub_service

import (
	"context"
	"fmt"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/sub_shard"
)

type ApiService struct {
	cm cluster.Manager
	UnimplementedSubServiceServer
}

func RegisterApiService(cs *service.ClusterService) {
	svc := &ApiService{
		cm: cs.ClusterManager(),
	}
	RegisterSubServiceServer(cs.GrpcServiceRegistrar(), svc)
	cs.AddGrpcGatewayRegister(func(ctx context.Context, mux *gruntime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		return RegisterSubServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", cs.Config().GrpcPort), opts)
	})

}

func (it *ApiService) client(shardID uint64) sub_shard.ISubShardDragonboatClient {
	dc := runtime.NewDragonboatClient(it.cm.NodeHost(), shardID)
	return sub_shard.NewSubShardDragonboatClient(dc)
}

func (it *ApiService) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetShard(ctx, &master_shard.GetShardRequest{Name: req.ShardName})
	if err != nil {
		return &GetResponse{}, fmt.Errorf("GetShard(%v) err: %w", req, err)
	}
	resp1, err := it.client(resp.Shard.ShardId).Get(ctx, &sub_shard.GetRequest{})
	if err != nil {
		return &GetResponse{}, fmt.Errorf("client.Get() err: %w", err)
	}
	return &GetResponse{Value: resp1.Value}, nil
}

func (it *ApiService) Incr(ctx context.Context, req *IncrRequest) (*IncrResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetShard(ctx, &master_shard.GetShardRequest{Name: req.ShardName})
	if err != nil {
		return &IncrResponse{}, fmt.Errorf("GetShard(%v) err: %w", req, err)
	}
	resp1, err := it.client(resp.Shard.ShardId).Incr(ctx, &sub_shard.IncrRequest{})
	return &IncrResponse{Value: resp1.Value}, nil
}
