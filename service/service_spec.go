package service

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"time"
	"errors"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/lo"
	"github.com/lni/dragonboat/v3"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
)

type SpecService struct {
	cm cluster.Manager
	client IClient
	UnimplementedSpecServer
}

func RegisterSpecService(cs *ClusterService, client IClient) {
	svc := &SpecService{cm: cs.ClusterManager(), client: client}
	RegisterSpecServer(cs.GrpcServiceRegistrar(), svc)
	cs.AddGrpcGatewayRegister(func(ctx context.Context, mux *gruntime.ServeMux) error {
		return RegisterSpecHandlerClient(ctx, mux, client.Spec())
	})
}

func (it *SpecService) AddShardSpec(ctx context.Context, req *AddShardSpecRequest) (*AddShardSpecResponse, error) {
	nodes := it.getNodeCreateViews()
	req2 := &master_shard.CreateShardRequest{Name: req.ShardName, Nodes: nodes, Replica: req.Replica}
	resp, err := it.cm.Client().MasterShard().CreateShard(ctx, req2)
	if err != nil {
		return &AddShardSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("Client().MasterShard().CreateShard(%s) err: %w", req.ShardName, err))
	}
	return &AddShardSpecResponse{Shard: ShardSpecSiri(resp.Shard)}, nil
}

func (it *SpecService) RemoveShardSpec(ctx context.Context, req *RemoveShardSpecRequest) (*RemoveShardSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().DeleteShard(ctx, &master_shard.DeleteShardRequest{Name: req.ShardName})
	if err != nil {
		return &RemoveShardSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("Client().MasterShard().DeleteShard(%s) err: %w", req.ShardName, err))
	}
	return &RemoveShardSpecResponse{Updated: resp.Updated}, nil
}

func (it *SpecService) RebalanceShardSpec(ctx context.Context, req *RebalanceShardSpecRequest) (*RebalanceShardSpecResponse, error) {
	nodes := it.getNodeCreateViews()
	req2 := &master_shard.UpdateShardRequest{Name: req.ShardName, Nodes: nodes, Replica: req.Replica}
	resp, err := it.cm.Client().MasterShard().UpdateShard(ctx, req2)
	if err != nil {
		return &RebalanceShardSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("Client().MasterShard().UpdateShard(%s) err: %w", req.ShardName, err))
	}
	return &RebalanceShardSpecResponse{Updated: resp.Updated}, nil
}

func (it *SpecService) getNodeCreateViews() []*master_shard.NodeCreateView {
	nodes := make([]*master_shard.NodeCreateView, 0, it.cm.Members().Nums())
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		nodes = append(nodes, &master_shard.NodeCreateView{NodeHostId: m.Meta.NodeHostId, Addr: m.Meta.RaftAddress})
		return true
	})
	return nodes
}

func (it *SpecService) ListShardSpec(ctx context.Context, req *ListShardSpecRequest) (*ListShardSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().ListShards(ctx, &master_shard.ListShardsRequest{})
	if err != nil {
		return &ListShardSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("Client().MasterShard().ListShards() err: %w", err))
	}
	specs := lo.Map(resp.Shards, WithNopIndex(ShardSpecSiri))

	var views []*ShardView
	if req.Status {
		views = it.getShardView(specs...)
	} else {
		views = lo.Map(specs, func(spec *ShardSpec, _ int) *ShardView {
			return &ShardView{Spec: spec}
		})
	}
	return &ListShardSpecResponse{Shards: views, Version: resp.StateVersion}, nil
}

func (it *SpecService) getShardStatus() map[uint64]*ShardStatus {
	result := make(map[uint64]*ShardStatus, 0)
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		if m.State != nil {
			for shardID, member := range m.State.Shards {
				view := &NodeStatusView{
					NodeHostId: m.Meta.NodeHostId,
					NodeId:     member.NodeId,
					Addr:       m.Meta.RaftAddress,
					Nodes:      member.Nodes,
					Leading:    member.IsLeader,
				}
				if status, ok := result[shardID]; ok {
					status.Nodes = append(status.Nodes, view)
				} else {
					result[shardID] = &ShardStatus{Nodes: []*NodeStatusView{view}}
				}
			}
			for shardID, nodeID := range m.State.LogShards {
				view := &LogNodeStatusView{
					NodeHostId: m.Meta.NodeHostId,
					NodeId:     nodeID,
					Addr:       m.Meta.RaftAddress,
				}
				if status, ok := result[shardID]; ok {
					status.LogNodes = append(status.LogNodes, view)
				} else {
					result[shardID] = &ShardStatus{LogNodes: []*LogNodeStatusView{view}}
				}
			}
		}
		return true
	})
	for _, status := range result {
		sort.Slice(status.Nodes, func(i int, j int) bool {
			return status.Nodes[i].NodeId < status.Nodes[j].NodeId
		})
		sort.Slice(status.LogNodes, func(i int, j int) bool {
			return status.LogNodes[i].NodeId < status.LogNodes[j].NodeId
		})
	}
	return result
}

func (it *SpecService) getShardView(specs ...*ShardSpec) []*ShardView {
	result := make([]*ShardView, 0, len(specs))
	status := it.getShardStatus()
	for _, spec := range specs {
		view := ShardView{Spec: spec}
		if s, ok := status[spec.ShardId]; ok {
			s.Type, s.TypeDetail = it.getShardStatusType(spec, s)
			view.Status = s
		}
		result = append(result, &view)
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Spec.ShardId < result[j].Spec.ShardId
	})
	return result
}

func (it *SpecService) getShardStatusType(spec *ShardSpec, status *ShardStatus) (ShardStatusType, *string) {
	specDict := lo.SliceToMap(spec.Nodes, func(v *NodeView) (uint64, string) {
		return v.NodeId, v.Addr
	})
	statusDict := lo.SliceToMap(status.Nodes, func(v *NodeStatusView) (uint64, string) {
		return v.NodeId, v.Addr
	})
	logDict := lo.SliceToMap(status.LogNodes, func(v *LogNodeStatusView) (uint64, string) {
		return v.NodeId, v.Addr
	})
	if !reflect.DeepEqual(specDict, statusDict) {
		return ShardStatusType_Updating, lo.ToPtr("status.Nodes not match with Spec")
	}
	nodes := lo.Reduce(status.Nodes, func(nodes map[uint64]string, view *NodeStatusView, index int) map[uint64]string {
		if reflect.DeepEqual(nodes, view.Nodes) {
			return nodes
		}
		return nil
	}, statusDict)
	if nodes == nil {
		return ShardStatusType_Updating, lo.ToPtr("status.Nodes not match with status.Nodes.Nodes")
	}
	if !reflect.DeepEqual(specDict, logDict) {
		return ShardStatusType_Updating, lo.ToPtr("status.LogNodes not match with Spec")
	}
	return ShardStatusType_Ready, nil
}

func (it *SpecService) GetShardSpec(ctx context.Context, req *GetShardSpecRequest) (*GetShardSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetShard(ctx, &master_shard.GetShardRequest{Name: req.ShardName})
	if err != nil {
		return &GetShardSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("Client().MasterShard().GetShard() err: %w", err))
	}
	views := it.getShardView(ShardSpecSiri(resp.Shard))
	return &GetShardSpecResponse{Shard: views[0]}, nil

}

func (it *SpecService) CreateMigration(ctx context.Context, req *CreateMigrationRequest) (*CreateMigrationResponse, error) {
	err := it.cm.Client().Migration(req.ShardId).Migrate(ctx, &runtime.DragonboatPrepareMigrationRequest{
		CompressionType: runtime.CompressionType_Snappy,
		Expire:          0,
	})
	if err != nil {
		return &CreateMigrationResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.Migration(%d).Migrate err: %w", req.ShardId, err))
	}
	return &CreateMigrationResponse{}, nil
}

func (it *SpecService) GetMigration(ctx context.Context, req *QueryMigrationRequest) (*QueryMigrationResponse, error) {
	resp, err := it.cm.Client().Migration(req.ShardId).QueryMigration(ctx, &runtime.DragonboatQueryMigrationRequest{})
	if err != nil {
		return &QueryMigrationResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.Migration(%d).QueryMigration err: %w", req.ShardId, err))
	}
	return &QueryMigrationResponse{State: MigrationStateSiri(resp.View, req.ShardId)}, nil
}

func (it *SpecService) ListMigration(ctx context.Context, req *ListMigrationRequest) (*ListMigrationResponse, error) {
	shards := make(map[uint64]struct{})
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		if m.State != nil {
			for shardID := range m.State.Shards {
				shards[shardID] = struct{}{}
			}
		}
		return true
	})
	migrations := make([]*MigrationStateListItem, 0, len(shards))
	for shardID := range shards {
		var item MigrationStateListItem
		resp, err := it.cm.Client().Migration(shardID).QueryMigration(ctx, &runtime.DragonboatQueryMigrationRequest{}, runtime.WithClientTimeout(time.Second))
		if err != nil {
			if errors.Is(err, dragonboat.ErrClusterNotFound) {
				if resp2, err := it.client.Spec().GetMigration(ctx, &QueryMigrationRequest{ShardId: shardID}); err != nil {
					item.Item = &MigrationStateListItem_Err{Err: MigrationStateErrorSiri(err, shardID)}
				} else {
					item.Item = &MigrationStateListItem_Data{Data: resp2.State}
				}
			} else {
				item.Item = &MigrationStateListItem_Err{Err: MigrationStateErrorSiri(err, shardID)}
			}
		} else {
			item.Item = &MigrationStateListItem_Data{Data: MigrationStateSiri(resp.View, shardID)}
		}
		migrations = append(migrations, &item)
	}
	sort.Slice(migrations, func(i int, j int) bool {
		return GetMigrationStateListItemShardID(migrations[i]) < GetMigrationStateListItemShardID(migrations[j])
	})
	return &ListMigrationResponse{Migrations: migrations}, nil
}

func (it *SpecService) UpdateNodeSpec(ctx context.Context, req *UpdateNodeSpecRequest) (*UpdateNodeSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().UpdateNode(ctx, &master_shard.UpdateNodeRequest{Node: req.Node})
	if err != nil {
		return &UpdateNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().UpdateNode(%s) err: %w", req.Node.NodeHostId, err))
	}
	return &UpdateNodeSpecResponse{Updated: resp.Updated}, nil
}
func (it *SpecService) GetNodeSpec(ctx context.Context, req *GetNodeSpecRequest) (*GetNodeSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetNode(ctx, &master_shard.GetNodeRequest{NodeHostId: req.NodeHostId})
	if err != nil {
		return &GetNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().GetNode(%s) err: %w", req.NodeHostId, err))
	}
	return &GetNodeSpecResponse{Node: resp.Node}, nil
}
func (it *SpecService) CordonNodeSpec(ctx context.Context, req *CordonNodeSpecRequest) (*CordonNodeSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetNode(ctx, &master_shard.GetNodeRequest{NodeHostId: req.NodeHostId})
	if err != nil {
		return &CordonNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().GetNode(%s) err: %w", req.NodeHostId, err))
	}
	if v, ok := resp.Node.Labels["builtin.exclude"]; ok && v == "true" {
		return &CordonNodeSpecResponse{Updated: 0}, nil
	}
	resp.Node.Labels["builtin.exclude"] = "true"
	resp2, err := it.cm.Client().MasterShard().UpdateNode(ctx, &master_shard.UpdateNodeRequest{Node: resp.Node})
	if err != nil {
		return &CordonNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().UpdateNode(%s) err: %w", req.NodeHostId, err))
	}
	return &CordonNodeSpecResponse{Updated: resp2.Updated}, nil
}
func (it *SpecService) UncordonNodeSpec(ctx context.Context, req *UncordonNodeSpecRequest) (*UncordonNodeSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().GetNode(ctx, &master_shard.GetNodeRequest{NodeHostId: req.NodeHostId})
	if err != nil {
		return &UncordonNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().GetNode(%s) err: %w", req.NodeHostId, err))
	}
	if _, ok := resp.Node.Labels["builtin.exclude"]; !ok {
		return &UncordonNodeSpecResponse{Updated: 0}, nil
	}
	delete(resp.Node.Labels, "builtin.exclude")
	resp2, err := it.cm.Client().MasterShard().UpdateNode(ctx, &master_shard.UpdateNodeRequest{Node: resp.Node})
	if err != nil {
		return &UncordonNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().UpdateNode(%s) err: %w", req.NodeHostId, err))
	}
	return &UncordonNodeSpecResponse{Updated: resp2.Updated}, nil
}
func (it *SpecService) DrainNodeSpec(ctx context.Context, req *DrainNodeSpecRequest) (*DrainNodeSpecResponse, error) {
	resp, err := it.cm.Client().MasterShard().ListShards(ctx, &master_shard.ListShardsRequest{NodeHostId: &req.NodeHostId})
	if err != nil {
		return &DrainNodeSpecResponse{}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().ListShards(%s) err: %w", req.NodeHostId, err))
	}
	nodes := it.getNodeCreateViews()
	var updated uint64
	for _, shard := range resp.Shards {
		resp2, err := it.cm.Client().MasterShard().UpdateShard(ctx, &master_shard.UpdateShardRequest{Name: shard.ShardName, Nodes: nodes})
		if err != nil {
			return &DrainNodeSpecResponse{Updated: updated}, DragonboatErrorToGrpcError(fmt.Errorf("client.MasterShard().UpdateShard(%s) err: %w", shard.ShardName, err))
		}
		updated += uint64(resp2.Updated)
	}
	return &DrainNodeSpecResponse{Updated: updated}, nil
}
