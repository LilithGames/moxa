package cluster

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/goutils/syncutil"

	"github.com/LilithGames/moxa/master_shard"
)

func OnOnceMasterShardReady(stopper *syncutil.Stopper, cm Manager) chan struct{} {
	ch := make(chan struct{}, 0)
	stopper.RunWorker(func() {
		duration := time.Second
		timer := time.NewTimer(duration)
		defer timer.Stop()
		for {
			if _, err := cm.Client().MasterShard().Healthz(context.TODO(), &master_shard.HealthzRequest{}, runtime.WithClientTimeout(time.Second)); err == nil {
				close(ch)
				return
			}
			ResetTimer(timer, duration)
			select {
			case <-timer.C:
			case <-stopper.ShouldStop():
				return
			}
		}
	})
	return ch
}

func UpdateNodeHostInfoWorker(cm Manager) Worker {
	return func(stopper *syncutil.Stopper) error {
		stopper.RunWorker(func() {
			memberChanged, sub1 := cm.EventBus().Subscribe(EventTopic_NodeHostMembershipChanged.String())
			defer sub1.Close()
			nodeReady, sub2 := cm.EventBus().Subscribe(EventTopic_NodeHostNodeReady.String())
			defer sub2.Close()
			leaderUpdated, sub3 := cm.EventBus().Subscribe(EventTopic_NodeHostLeaderUpdated.String())
			defer sub3.Close()

			duration := time.Second * 30
			timer := time.NewTimer(duration)
			defer timer.Stop()
			for {
				nhi := cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
				if err := cm.Members().UpdateNodeHostInfo(nhi); err != nil {
					log.Println("[WARN]", fmt.Errorf("Members().UpdateNodeHostInfo err: %w", err))
				}
				if err := cm.Members().SyncState(); err != nil {
					log.Println("[WARN]", fmt.Errorf("Members().SyncState err: %w", err))
					continue
				}
				ResetTimer(timer, duration)
				select {
				case <-memberChanged:
				case <-nodeReady:
				case <-leaderUpdated:
				case <-timer.C:
				case <-stopper.ShouldStop():
					return
				}
			}

		})
		return nil
	}
}
