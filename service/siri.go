package service

import (
	"fmt"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/samber/lo"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"

	"github.com/LilithGames/moxa/master_shard"
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

func LogShardInfoSiri(ni *raftio.NodeInfo) *LogShardInfo {
	return &LogShardInfo{
		ShardId: ni.ClusterID,
		NodeId: ni.NodeID,
	}
}

func NodeHostInfoSiri(src *dragonboat.NodeHostInfo) *NodeHostInfo {
	return &NodeHostInfo{
		NodeHostId: src.NodeHostID,
		RaftAddress: src.RaftAddress,
		ShardInfoList: lo.Map(src.ClusterInfoList, func(src dragonboat.ClusterInfo, _ int) *ShardInfo {
			return ShardInfoSiri(&src)
		}),
		LogShardInfoList: lo.Map(src.LogInfo, func(ni raftio.NodeInfo, _ int) *LogShardInfo {
			return LogShardInfoSiri(&ni)
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
		Initials: lo.MapToSlice(src.Initials, func(_ string, node *master_shard.Node) *NodeView {
			return NodeViewSiri(node)
		}),
		Nodes: lo.MapToSlice(src.Nodes, func(_ string, node *master_shard.Node) *NodeView {
			return NodeViewSiri(node)
		}),
		Replica: src.Replica,
		LastNodeId: src.LastNodeId,
	}
}

func MigrationNodeSiri(src *runtime.MigrationNodeWrapper) *MigrationNodeView {
	return &MigrationNodeView{
		NodeId: src.Node.NodeId,
		Version: src.Node.Version,
		Nodes: lo.MapToSlice(src.Node.Nodes, func(nid uint64, _ bool) uint64 {
			return nid
		}),
		Epoch: src.Epoch,
	}
}

func MigrationStateSiri(src *runtime.MigrationStateView, shardID uint64) *MigrationState {
	return &MigrationState{
		ShardId: shardID,
		Type: src.Type.String(),
		Version: src.Version,
		Epoch: src.Epoch,
		NodeUpgraded: lo.MapValues(src.NodeUpgraded, func(v *runtime.MigrationNodeWrapper, _ uint64) *MigrationNodeView {
			return MigrationNodeSiri(v)
		}),
		NodeSaved: lo.MapValues(src.NodeSaved, func(v *runtime.MigrationNodeWrapper, _ uint64) *MigrationNodeView {
			return MigrationNodeSiri(v)
		}),            
		StateIndex: src.StateIndex,
	}
}
func MigrationStateErrorSiri(err error, shardID uint64) *MigrationStateError {
	return &MigrationStateError{ShardId: shardID, Error: err.Error()}
}

func GetMigrationStateListItemShardID(item *MigrationStateListItem) uint64 {
	data := item.GetData()
	if data != nil {
		return data.ShardId
	}
	err := item.GetErr()
	if err != nil {
		return err.ShardId
	}
	return 0
}

func WithNopIndex[T any, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R {
		return f(t)
	}
}

