package moxa

import (
	"context"
	"fmt"
	"log"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/goutils/syncutil"

	"github.com/LilithGames/moxa/cluster"
)

func ShardSpecMembershipWorker(cm cluster.Manager, sm *ShardManager) cluster.Worker {
	return func(stopper *syncutil.Stopper) error {
		handle := func(ev eventbus.Event) {
			defer ev.Done()
			e := ev.Data.(*cluster.ShardSpecChangingEvent)
			ctx := context.TODO()
			switch e.Type {
			case cluster.ShardSpecChangingType_Adding:
				if err := sm.ServeSubShard(ctx, e.ShardId, *e.CurrentNodeId, e.CurrentNodes); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_Adding ServeSubShard(%d) err: %w", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d added", e.ShardId))
				}
			case cluster.ShardSpecChangingType_Deleting:
				if err := sm.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, false); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_Deleting StopSubShard(%d) err: %w", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d deleted", e.ShardId))
				}
			case cluster.ShardSpecChangingType_NodeJoining:
				if err := sm.ServeSubShard(ctx, e.ShardId, *e.CurrentNodeId, map[uint64]string{}); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_NodeJoining ServeSubShard(%d) err: %v", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d joined", e.ShardId))
				}
			case cluster.ShardSpecChangingType_NodeLeaving:
				if err := sm.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, true); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_NodeLeaving StopSubShard(%d) err: %v", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d leaved", e.ShardId))
				}
			case cluster.ShardSpecChangingType_NodeMoving:
				if err := sm.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, true); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_NodeMoving StopSubShard(%d) err: %v", e.ShardId, err))
				} else if err := sm.ServeSubShard(ctx, e.ShardId, *e.CurrentNodeId, map[uint64]string{}); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.SubShardChangingType_NodeMoving ServeSubShard(%d) err: %v", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d moved", e.ShardId))
				}
			case cluster.ShardSpecChangingType_Cleanup:
				if err := sm.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, false); err != nil {
					log.Println("[ERROR]", fmt.Errorf("ShardSpecMembershipWorker.ShardSpecChangingType_Cleanup StopSubShard(%d) err: %w", e.ShardId, err))
				} else {
					log.Println("[INFO]", fmt.Sprintf("ShardSpecMembershipWorker: SubShard %d cleanup success", e.ShardId))
				}
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
