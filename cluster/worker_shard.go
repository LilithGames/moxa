package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/goutils/syncutil"
	"github.com/samber/lo"
	"github.com/lni/dragonboat/v3/raftio"

	"github.com/LilithGames/moxa/master_shard"
)

func ShardSpecUpdatingWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		stopper.RunWorker(func() {
			ctx := context.TODO()
			tick := time.NewTicker(time.Second)
			defer tick.Stop()
			var version uint64 = 0
			for {
				select {
				case <-tick.C:
					resp, err := cm.Client().MasterShard().GetStateVersion(ctx, &master_shard.GetStateVersionRequest{}, runtime.WithClientStale(true))
					if err != nil {
						if !errors.Is(err, dragonboat.ErrClusterNotFound) && !errors.Is(err, dragonboat.ErrClosed) {
							log.Println("[WARN]", fmt.Errorf("OnSubShardsUpdating MasterShard().GetStateVersion err: %w", err))
						}
						continue
					}
					if version != resp.ShardsStateVersion {
						if version > resp.ShardsStateVersion {
							log.Println("[WARN]", fmt.Errorf("OnSubShardsUpdating ShardsStateVersion seems rollbacked"))
						} else {
							log.Println("[INFO]", fmt.Errorf("OnSubShardsUpdating shards version %d -> %d", version, resp.ShardsStateVersion))
						}
						cm.EventBus().Publish(EventTopic_ShardSpecUpdating.String(), &ShardSpecUpdatingEvent{StateVersion: resp.ShardsStateVersion})
						version = resp.ShardsStateVersion
					}
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}

func diffSubShards(addr string, spec []*master_shard.ShardSpec, status []dragonboat.ClusterInfo, persist []raftio.NodeInfo) []*ShardSpecChangingEvent {
	specDict := lo.SliceToMap(spec, func(s *master_shard.ShardSpec) (uint64, map[uint64]string) {
		return s.ShardId, lo.MapEntries(s.Nodes, func(nhid string, node *master_shard.Node) (uint64, string) {
			return node.NodeId, node.Addr
		})
	})
	statusDict := lo.SliceToMap(status, func(ci dragonboat.ClusterInfo) (uint64, map[uint64]string) {
		return ci.ClusterID, ci.Nodes
	})
	persistDict := lo.SliceToMap(persist, func(ni raftio.NodeInfo) (uint64, uint64) {
		return ni.ClusterID, ni.NodeID
	})
	specKeys := lo.MapToSlice(specDict, func(sid uint64, _ map[uint64]string) uint64 {
		return sid
	})
	statusKeys := lo.MapToSlice(statusDict, func(sid uint64, _ map[uint64]string) uint64 {
		return sid
	})
	persistKeys := lo.MapToSlice(persistDict, func(sid uint64, _ uint64) uint64 {
		return sid
	})
	keys := lo.Union(lo.Union(specKeys, statusKeys), persistKeys)
	results := lo.FilterMap(keys, func(sid uint64, _ int) (*ShardSpecChangingEvent, bool) {
		if sid == MasterShardID {
			return nil, false
		}
		specNode, specOK := specDict[sid]
		statusNode, statusOK := statusDict[sid]
		persistNodeID, persistOK := persistDict[sid]
		specNodeID, specNodeOK := lo.FindKey(specNode, addr)
		statusNodeID, statusNodeOK := lo.FindKey(statusNode, addr)
		event := &ShardSpecChangingEvent{ShardId: sid, PreviousNodes: statusNode, CurrentNodes: specNode}
		if specOK && !statusOK {
			if !specNodeOK {
				log.Println("[WARN] diffSubShards impossible cause: SubShardChangingType_Adding !specNodeOK")
				return nil, false
			}
			event.Type = ShardSpecChangingType_Adding
			event.CurrentNodeId = &specNodeID
			return event, true
		} else if !specOK && statusOK {
			if !statusNodeOK {
				log.Println("[WARN] diffSubShards impossible cause: SubShardChangingType_Deleting !statusNodeOK")
				return nil, false
			}
			event.Type = ShardSpecChangingType_Deleting
			event.PreviousNodeId = &statusNodeID
			return event, true
		} else if specOK && statusOK {
			if specNodeOK && !statusNodeOK {
				event.Type = ShardSpecChangingType_NodeJoining
				event.CurrentNodeId = &specNodeID
				return event, true
			} else if !specNodeOK && statusNodeOK {
				event.Type = ShardSpecChangingType_NodeLeaving
				event.PreviousNodeId = &statusNodeID
				return event, true
			} else if specNodeOK && statusNodeOK {
				if specNodeID != statusNodeID {
					event.Type = ShardSpecChangingType_NodeMoving
					event.CurrentNodeId = &specNodeID
					event.PreviousNodeId = &statusNodeID
					return event, true
				}
			}
			return nil, false
		} else if persistOK {
			event.Type = ShardSpecChangingType_Cleanup
			event.PreviousNodeId = &persistNodeID
			return event, true
		}
		return nil, false
	})
	return results
}

func ShardSpecChangingWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		handle := func(e *eventbus.Event) {
			if e != nil {
				defer e.Done()
			}
			resp, err := cm.Client().MasterShard().ListShards(context.TODO(), &master_shard.ListShardsRequest{NodeHostId: lo.ToPtr(cm.NodeHostID())}, runtime.WithClientStale(true))
			if err != nil {
				log.Println("[ERROR]", fmt.Errorf("ShardSpecChangingWorker MasterShard.ListShards err: %w", err))
				return
			}
			nhi := cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})
			events := diffSubShards(nhi.RaftAddress, resp.Shards, nhi.ClusterInfoList, nhi.LogInfo)
			log.Println("[INFO]", fmt.Sprintf("ShardSpecChangingWorker %d shard changes detected", len(events)))
			for _, event := range events {
				cm.EventBus().Publish(EventTopic_ShardSpecChanging.String(), event)
			}
		}
		stopper.RunWorker(func() {
			shardsUpdating, sub := cm.EventBus().Subscribe(EventTopic_ShardSpecUpdating.String())
			defer sub.Close()
			duration := time.Second * 30
			timer := time.NewTimer(duration)
			defer timer.Stop()
			for {
				ResetTimer(timer, duration)
				select {
				case <-timer.C:
					handle(nil)
				case e := <-shardsUpdating:
					handle(&e)
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}
