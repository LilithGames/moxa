package service

import (
	"fmt"

	"github.com/lni/dragonboat/v3"
	"github.com/samber/lo"
	sm "github.com/lni/dragonboat/v3/statemachine"

	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/cluster"
)

func StateMachineTypeSiri(src sm.Type) string {
	switch src {
	case sm.RegularStateMachine:
		return "RegularStateMachine"
	case sm.ConcurrentStateMachine:
		return "ConcurrentStateMachine"
	case sm.OnDiskStateMachine:
		return "OnDiskStateMachine"
	default:
		return fmt.Sprintf("%d", src)
	}
}

func ShardInfoSiri(src *dragonboat.ClusterInfo) *ShardInfo {
	return &ShardInfo{
		ShardId: src.ClusterID,
		NodeId: src.NodeID,
		Nodes:  src.Nodes,
		StateMachineType: StateMachineTypeSiri(src.StateMachineType),
		IsLeader: src.IsLeader,
		IsObserver: src.IsObserver,
		IsWitness: src.IsWitness,
		Pending: src.Pending,
	}
}

func NodeHostInfoSiri(src *dragonboat.NodeHostInfo) *NodeHostInfo {
	return &NodeHostInfo{
		NodeHostId: src.NodeHostID,
		RaftAddress: src.RaftAddress,
		ShardInfoList: lo.Map(src.ClusterInfoList, func(src dragonboat.ClusterInfo, _ int) *ShardInfo {
			return ShardInfoSiri(&src)
		}),
	}
}

func NodeViewSiri(src *master_shard.Node) *NodeView {
	return &NodeView{
		NodeHostId: src.NodeHostId,
		NodeId: src.NodeId,
		Addr: src.Addr,
	}
}

func ShardSpecSiri(src *master_shard.ShardSpec) *ShardSpec {
	return &ShardSpec{
		ShardName: src.ShardName,
		ShardId: src.ShardId,
		Nodes: lo.MapToSlice(src.Nodes, func(_ string, node *master_shard.Node) *NodeView {
			return NodeViewSiri(node)
		}),
	}
}

func CreateSnapshotResponseSiri(src *cluster.Snapshot, nodes map[uint64]string) *CreateSnapshotResponse {
	return &CreateSnapshotResponse{
		SnapshotIndex: src.Index,
		Version: src.Version,
		Snapshot: &RemoteSnapshot{
			Index: src.Index,
			Version: src.Version,
			Path: src.Path(),
			NodeHostId: src.NodeHostID,
			Nodes: nodes,
		},
	}
}

func WithNopIndex[T any, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R {
		return f(t)
	}
}
