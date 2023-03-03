package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/goutils/syncutil"

	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/utils"
)

func ShardSpecUpdatingWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		stopper.RunWorker(func() {
			ctx := context.TODO()
			tick := time.NewTicker(time.Millisecond*100)
			defer tick.Stop()
			var version uint64 = 0
			for {
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
				select {
				case <-tick.C:
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}

func ShardSpecChangingWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		maxDuration := time.Second * 30
		handle := func(e *eventbus.Event) time.Duration {
			if e != nil {
				defer e.Done()
			}
			events, err := ReconcileSubShardChangingEvents(cm)
			if err != nil {
				log.Println("[ERROR]", fmt.Errorf("ShardSpecChangingWorker GetSubShardChangingEvents err: %w", err))
				return maxDuration
			}
			if len(events) > 0 {
				log.Println("[INFO]", fmt.Sprintf("ShardSpecChangingWorker %d shard changes detected", len(events)))
			}
			for _, event := range events {
				cm.EventBus().Publish(EventTopic_ShardSpecChanging.String(), event)
			}
			if len(events) > 0 {
				return time.Second * 3
			}
			return maxDuration
		}
		stopper.RunWorker(func() {
			shardsUpdating, sub := cm.EventBus().Subscribe(EventTopic_ShardSpecUpdating.String())
			defer sub.Close()
			duration := maxDuration
			timer := time.NewTimer(duration)
			defer timer.Stop()
			for {
				utils.ResetTimer(timer, duration)
				select {
				case <-timer.C:
					duration = handle(nil)
				case e := <-shardsUpdating:
					duration = handle(&e)
				case <-stopper.ShouldStop():
					return
				}
			}
		})
		return nil
	}
}
