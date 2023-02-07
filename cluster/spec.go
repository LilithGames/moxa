package cluster

import (
	"context"
	"fmt"
	"log"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/utils"
)

func ReconcileSubShardChangingEvents(cm Manager) ([]*ShardSpecChangingEvent, error) {
	resp, err := cm.Client().MasterShard().ListShards(context.TODO(), &master_shard.ListShardsRequest{}, runtime.WithClientStale(true))
	if err != nil {
		return nil, fmt.Errorf("MasterShard().ListShards stale err: %w", err)
	}
	nhi := cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})
	events := DiffSubShards(nhi.RaftAddress, resp.Shards, nhi.ClusterInfoList, nhi.LogInfo)
	return events, nil
}

type ShardSpecDiffView struct {
	Spec *master_shard.ShardSpec
	Nodes map[uint64]string
	Initials map[uint64]string
}
type ShardStatusDiffView struct {
	Status dragonboat.ClusterInfo
	Nodes map[uint64]string
}

func DiffSubShards(addr string, spec []*master_shard.ShardSpec, status []dragonboat.ClusterInfo, persist []raftio.NodeInfo) []*ShardSpecChangingEvent {
	specDict := lo.SliceToMap(spec, func(s *master_shard.ShardSpec) (uint64, *ShardSpecDiffView) {
		return s.ShardId, &ShardSpecDiffView{
			Spec: s,
			Nodes: lo.MapEntries(s.Nodes, func(nhid string, node *master_shard.Node) (uint64, string) {
				return node.NodeId, node.Addr
			}),
			Initials: lo.MapEntries(s.Initials, func(nhid string, node *master_shard.Node) (uint64, string) {
				return node.NodeId, node.Addr
			}),
		}
	})
	statusDict := lo.SliceToMap(status, func(ci dragonboat.ClusterInfo) (uint64, *ShardStatusDiffView) {
		return ci.ClusterID, &ShardStatusDiffView{
			Status: ci,
			Nodes: ci.Nodes,
		}
	})
	persistDict := lo.SliceToMap(persist, func(ni raftio.NodeInfo) (uint64, uint64) {
		return ni.ClusterID, ni.NodeID
	})
	specKeys := lo.MapToSlice(specDict, func(sid uint64, _ *ShardSpecDiffView) uint64 {
		return sid
	})
	statusKeys := lo.MapToSlice(statusDict, func(sid uint64, _ *ShardStatusDiffView) uint64 {
		return sid
	})
	persistKeys := lo.MapToSlice(persistDict, func(sid uint64, _ uint64) uint64 {
		return sid
	})
	keys := lo.Union(lo.Union(specKeys, statusKeys), persistKeys)
	results := lo.FilterMap(keys, func(sid uint64, _ int) ([]*ShardSpecChangingEvent, bool) {
		if sid == MasterShardID {
			return nil, false
		}
		specView, specOK := specDict[sid]
		statusView, statusOK := statusDict[sid]
		persistNodeID, persistOK := persistDict[sid]
		event := &ShardSpecChangingEvent{ShardId: sid}
		if specOK {
			event.CurrentNodes = specView.Nodes
			event.ProfileName = &specView.Spec.ProfileName
		}
		if statusOK {
			event.PreviousNodes = statusView.Nodes
		}

		if specOK && !statusOK {
			specNodeID, specNodeOK := lo.FindKey(specView.Nodes, addr)
			if specNodeOK && persistOK {
				event.Type = ShardSpecChangingType_Recovering
				event.PreviousNodeId = &persistNodeID
				return utils.ToSlice(event), true
			} else if specNodeOK && !persistOK {
				if _, ok := specView.Initials[specNodeID]; ok {
					event.Type = ShardSpecChangingType_Adding
					event.CurrentNodeId = &specNodeID
					return utils.ToSlice(event), true
				} else {
					event.Type = ShardSpecChangingType_Joining
					event.CurrentNodeId = &specNodeID
					return utils.ToSlice(event), true
				}
			} else if !specNodeOK && persistOK {
				event.Type = ShardSpecChangingType_Cleanup
				event.PreviousNodeId = &persistNodeID
				return utils.ToSlice(event), true
			}
			return nil, false
		} else if !specOK && statusOK {
			statusNodeID, statusNodeOK := lo.FindKey(statusView.Nodes, addr)
			if !statusNodeOK {
				log.Println("[WARN] diffSubShards impossible cause: !specOK && statusOK !statusNodeOK")
				return nil, false
			}
			event.Type = ShardSpecChangingType_Deleting
			event.PreviousNodeId = &statusNodeID
			return utils.ToSlice(event), true
		} else if specOK && statusOK {
			statusNodeID, statusNodeOK := lo.FindKey(statusView.Nodes, addr)
			specNodeID, specNodeOK := lo.FindKey(specView.Nodes, addr)
			if !statusNodeOK {
				log.Println("[WARN] diffSubShards impossible cause: specOK && statusOK but !statusNodeOK")
				return nil, false
			}
			if !specNodeOK || specNodeID != statusNodeID {
				event.Type = ShardSpecChangingType_Leaving
				event.PreviousNodeId = &statusNodeID
				return utils.ToSlice(event), true
			} else {
				// join or leave other nodes
				events := DiffNodes(specView.Nodes, statusView.Nodes, event)
				return events, true
			}
		} else if persistOK {
			event.Type = ShardSpecChangingType_Cleanup
			event.PreviousNodeId = &persistNodeID
			return utils.ToSlice(event), true
		}
		return nil, false
	})
	return lo.Flatten(results)
}

func DiffNodes(spec map[uint64]string, status map[uint64]string, tmpl *ShardSpecChangingEvent) []*ShardSpecChangingEvent {
	specKeys := lo.MapToSlice(spec, func(sid uint64, _ string) uint64 {
		return sid
	})
	statusKeys := lo.MapToSlice(status, func(sid uint64, _ string) uint64 {
		return sid
	})
	keys := lo.Union(specKeys, statusKeys)
	events := lo.FilterMap(keys, func(nodeID uint64, _ int) (*ShardSpecChangingEvent, bool) {
		specAddr, specOK := spec[nodeID]
		statusAddr, statusOK := status[nodeID]
		event := utils.Clone(tmpl)
		if specOK && statusOK {
			if specAddr != statusAddr {
				log.Println("[WARN]", fmt.Errorf("DiffNodes: impossible case specOK && statusOK but specAddr != statusAddr"))
			}
			return nil, false
		} else if specOK && !statusOK {
			event.Type = ShardSpecChangingType_NodeJoining
			event.CurrentNodeId = &nodeID
			return event, true
		} else if !specOK && statusOK {
			event.Type = ShardSpecChangingType_NodeLeaving
			event.PreviousNodeId = &nodeID
			return event, true
		}
		return nil, false
	})
	return events
}
