package cluster

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lni/goutils/syncutil"
	"github.com/puzpuzpuz/xsync"
	"github.com/samber/lo"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"

	"github.com/LilithGames/moxa/utils"
)

func MigrationChangedWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		stopper.RunWorker(func() {
			ctx := context.TODO()
			index := xsync.NewTypedMapOf[uint64, uint64](func(k uint64) uint64 { return k })
			handle := func(force bool) {
				shardIDs := lo.MapToSlice(cm.GetNodeIDs(), func(sid uint64, _ uint64) uint64 {
					return sid
				})
				_, err := utils.ParallelMap(shardIDs, func(shardID uint64) (uint64, error) {
					resp, err := cm.Client().Migration(shardID).QueryIndex(ctx, &runtime.DragonboatQueryIndexRequest{}, runtime.WithClientStale(true))
					if err != nil {
						return 0, fmt.Errorf("Client().Migration(%d).QueryIndex() err: %w", shardID, err)
					}
					if prev, ok := index.Load(shardID); !ok {
						cm.EventBus().Publish(EventTopic_MigrationChanged.String(), &MigrationChangedEvent{ShardId: shardID, StateIndex: resp.StateIndex})
						index.Store(shardID, resp.StateIndex)
					} else if resp.StateIndex > prev || force {
						cm.EventBus().Publish(EventTopic_MigrationChanged.String(), &MigrationChangedEvent{ShardId: shardID, StateIndex: resp.StateIndex})
						index.Store(shardID, resp.StateIndex)
					}
					return resp.StateIndex, nil
				})
				if err != nil {
					log.Println("[WARN]", fmt.Errorf("MigrationUpdating ParallelMap err: %w", err))
				}
			}

			duration := time.Second * 30
			timer := time.NewTimer(duration)
			defer timer.Stop()
			tick := time.NewTicker(time.Minute)
			defer tick.Stop()
			memberChanged, sub := cm.EventBus().Subscribe(EventTopic_NodeHostMembershipChanged.String())
			defer sub.Close()

			handle(false)
			utils.ResetTimer(timer, duration)
			for {
				select {
				case <-timer.C:
					handle(false)
					utils.ResetTimer(timer, duration)
				case <-memberChanged:
					handle(false)
				case <-tick.C:
					handle(true)
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}
