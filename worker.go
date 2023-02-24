package moxa

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/goutils/syncutil"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/utils"
)

func MigrationNodeWorker(cm cluster.Manager, sm *ShardManager) cluster.Worker {
	return func(stopper *syncutil.Stopper) error {
		ctx := utils.BindContext(stopper, context.TODO())
		stopper.RunWorker(func() {
			changed, sub := cm.EventBus().Subscribe(cluster.EventTopic_MigrationChanged.String())
			defer sub.Close()
			handle := func(ev eventbus.Event) {
				defer ev.Done()
				e := ev.Data.(*cluster.MigrationChangedEvent)
				if err := sm.ReconcileMigration(ctx, e.ShardId); err != nil {
					log.Println("[WARN]", fmt.Errorf("MigrationNodeWorker ReconcileMigration(%d) err: %w", e.ShardId, err))
				}
			}
			for {
				select {
				case ev := <-changed:
					stopper.RunWorker(func() {
						handle(ev)
					})
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}

func ShardSpecMembershipWorker(cm cluster.Manager, sm *ShardManager) cluster.Worker {
	return func(stopper *syncutil.Stopper) error {
		handle := func(ev eventbus.Event) {
			defer ev.Done()
			e := ev.Data.(*cluster.ShardSpecChangingEvent)
			ctx := context.TODO()
			if err := sm.ShardSpecChangingDispatcher(ctx, e); err != nil {
				log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker shard %d %s err: %w", e.ShardId, e.Type.String(), err))
			} else {
				var nodeID uint64
				if e.PreviousNodeId != nil {
					nodeID = *e.PreviousNodeId
				} else if e.CurrentNodeId != nil {
					nodeID = *e.CurrentNodeId
				}
				log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker shard %d %s node %d success", e.ShardId, e.Type.String(), nodeID))
			}
		}
		stopper.RunWorker(func() {
			shardsChanging, sub := cm.EventBus().Subscribe(cluster.EventTopic_ShardSpecChanging.String())
			defer sub.Close()
			for {
				select {
				case ev := <-shardsChanging:
					handle(ev)
				case <-stopper.ShouldStop():
					return
				}
			}

		})
		return nil
	}
}

func ShardSpecFixedSizeCreationWorker(cm cluster.Manager, profileName string, size int) cluster.Worker {
	return func(stopper *syncutil.Stopper) error {
		ctx := context.TODO()
		parse := func(name string) *uint64 {
			items := strings.Split(name, "-")
			if len(items) != 2 {
				return nil
			}
			if items[0] != "auto" {
				return nil
			}
			result, err := strconv.ParseUint(items[1], 10, 64)
			if err != nil {
				return nil
			}
			return &result
		}
		maxID := func(specs []*master_shard.ShardSpec) uint64 {
			ids := lo.Map(specs, func(spec *master_shard.ShardSpec, _ int) uint64 {
				i := parse(spec.ShardName)
				if i == nil {
					return 0
				}
				return *i
			})
			return lo.Max(ids)
		}
		bo := &utils.Backoff{
			InitialDelay: time.Millisecond,
			Step:         2,
			MaxDuration:  lo.ToPtr(time.Second*10),
		}
		reconcile := func() bool {
			resp, err := cm.Client().MasterShard().ListShards(ctx, &master_shard.ListShardsRequest{}, runtime.WithClientTimeout(time.Second*3))
			if err != nil {
				log.Println("[WARN]", fmt.Errorf("ShardSpecFixedSizeCreationWorker MasterShard().ListShards err: %w", err))
				return false
			}
			count := size - len(resp.Shards)
			start := time.Now()
			if count > 0 {
				id := maxID(resp.Shards)
				id += 1
				nodes := make([]*master_shard.NodeCreateView, 0, cm.Members().Nums())
				cm.Members().Foreach(func(m cluster.MemberNode) bool {
					nodes = append(nodes, &master_shard.NodeCreateView{NodeHostId: m.Meta.NodeHostId, Addr: m.Meta.RaftAddress})
					return true
				})
				name := fmt.Sprintf("auto-%d", id)
				nodesView := &master_shard.ShardNodesView{Nodes: nodes, Replica: lo.ToPtr(cm.MinClusterSize())}
				resp, err := cm.Client().MasterShard().CreateShard(ctx, &master_shard.CreateShardRequest{Name: name, ProfileName: profileName, NodesView: nodesView}, runtime.WithClientTimeout(time.Second*10))
				if err != nil {
					log.Println("[WARN]", fmt.Errorf("ShardSpecFixedSizeCreationWorker MasterShard().CreateShard %s err: %w", name, err))
					return false
				}
				if err := utils.RetryWithBackoff(ctx, bo, func() (bool, error) {
					if err := cm.Client().Void(resp.Shard.ShardId).VoidQuery(ctx, runtime.WithClientTimeout(time.Second)); err != nil {
						return true, fmt.Errorf("VoidQuery err: %w", err)
					}
					return false, nil
				}); err != nil {
					log.Println("[WARN]", fmt.Errorf("ShardSpecFixedSizeCreationWorker %s healthz err: %w", resp.Shard.ShardName, err))
					return false
				}
				log.Println("[INFO]", fmt.Sprintf("ShardSpecFixedSizeCreationWorker create shard %s success by %s", resp.Shard.ShardName, time.Since(start)))
			} else if count < 0 {
				last := resp.Shards[len(resp.Shards)-1]
				if _, err := cm.Client().MasterShard().DeleteShard(ctx, &master_shard.DeleteShardRequest{Name: last.ShardName}, runtime.WithClientTimeout(time.Second*10)); err != nil {
					log.Println("[WARN]", fmt.Errorf("ShardSpecFixedSizeCreationWorker DeleteShard %s err: %w", last.ShardName, err))
					return false
				}
				if err := utils.RetryWithBackoff(ctx, bo, func() (bool, error) {
					if err := cm.Client().Void(last.ShardId).VoidQuery(ctx, runtime.WithClientTimeout(time.Second)); err == nil {
						return true, fmt.Errorf("shard %d still alive", last.ShardId)
					}
					return false, nil
				}); err != nil {
					log.Println("[WARN]", fmt.Errorf("ShardSpecFixedSizeCreationWorker delete failed"))
					return false
				}
				log.Println("[INFO]", fmt.Sprintf("ShardSpecFixedSizeCreationWorker delete shard %s success by %s", last.ShardName, time.Since(start)))
			}
			return true
		}
		stopper.RunWorker(func() {
			duration := time.Second * 15
			timer := time.NewTimer(duration)
			defer timer.Stop()
			for {
				if !cm.IsLeader(cluster.MasterShardID) {
					duration = time.Second * 15
				} else if reconcile() {
					duration = time.Millisecond * 100
				} else {
					duration = time.Minute
				}
				utils.ResetTimer(timer, duration)
				select {
				case <-timer.C:
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}
