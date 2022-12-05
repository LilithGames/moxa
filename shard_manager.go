package moxa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/hashicorp/go-multierror"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/lni/goutils/syncutil"
	"github.com/samber/lo"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/utils"
)

var ErrInitialMemberNotEnough error = errors.New("ErrInitialMemberNotEnough")

type ShardManager struct {
	cm      cluster.Manager
	client  service.IClient
	master  any
	fn      any
	stopper *syncutil.Stopper
}

func NewShardManager(cm cluster.Manager, client service.IClient, master any, fn any) (*ShardManager, error) {
	return &ShardManager{cm, client, master, fn, syncutil.NewStopper()}, nil
}

func (it *ShardManager) Stop() {
	it.stopper.Stop()
}

func (it *ShardManager) getMasterInitialMembers() map[uint64]string {
	imembers := make(map[uint64]string, 0)
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		if m.Meta.NodeHostIndex < it.cm.MinClusterSize() {
			imembers[m.Meta.MasterNodeId] = m.Meta.RaftAddress
		}
		return true
	})
	return imembers
}

func (it *ShardManager) ServeMasterShard(ctx context.Context) error {
	shardID := cluster.MasterShardID
	conf := it.GetShardConfig(shardID, it.cm.MasterNodeID())
	if it.cm.NodeHost().HasNodeInfo(shardID, it.cm.MasterNodeID()) {
		// recovery
		if it.cm.GetNodeID(cluster.MasterShardID) != nil {
			// already recovery, skip
			return nil
		}
		log.Printf("[INFO] shardmanager: decided to recovery shard %d\n", shardID)
		if err := it.StartCluster(map[uint64]string{}, false, it.master, conf); err != nil {
			if errors.Is(err, dragonboat.ErrClusterAlreadyExist) {
				return nil
			}
			return fmt.Errorf("RecoveryShard %d err: %w", shardID, err)
		}
	} else {
		if it.cm.NodeHostIndex() < it.cm.MinClusterSize() {
			// create shard
			log.Printf("[INFO] shardmanager: decided to create shard %d\n", shardID)
			initials := it.getMasterInitialMembers()
			if len(initials) < int(it.cm.MinClusterSize()) {
				return ErrInitialMemberNotEnough
			}
			if err := it.StartCluster(initials, false, it.master, conf); err != nil {
				return fmt.Errorf("CreateShard %d err: %w", shardID, err)
			}
		} else {
			// join shard
			log.Printf("[INFO] shardmanager: decided to join shard %d\n", shardID)
			if err := it.JoinShard(ctx, shardID, it.cm.MasterNodeID()); err != nil {
				return fmt.Errorf("JoinShard err: %w", err)
			}
			if err := it.StartCluster(map[uint64]string{}, true, it.master, conf); err != nil {
				return fmt.Errorf("JoinShard %d err: %w", shardID, err)
			}
		}
	}
	return nil
}

func (it *ShardManager) ServeSubShard(ctx context.Context, shardID uint64, nodeID uint64, initials map[uint64]string) error {
	if shardID == cluster.MasterShardID {
		return fmt.Errorf("don't use ServeSubShard to serve master-shard")
	}
	if nodeID == 0 {
		return fmt.Errorf("nodeID must not be zero")
	}
	conf := it.GetShardConfig(shardID, nodeID)
	if it.cm.NodeHost().HasNodeInfo(shardID, nodeID) {
		// recovery
		log.Printf("[INFO] shardmanager: decided to recovery subshard %d, nodeID: %d\n", shardID, nodeID)
		if err := it.StartCluster(map[uint64]string{}, false, it.fn, conf); err != nil {
			return fmt.Errorf("RecoverySubShard %d err: %w", shardID, err)
		}
	} else {
		if len(initials) > 0 {
			// create shard
			log.Printf("[INFO] shardmanager: decided to create subshard %d, nodeID: %d\n", shardID, nodeID)
			if err := it.StartCluster(initials, false, it.fn, conf); err != nil {
				return fmt.Errorf("CreateSubShard %d err: %w", shardID, err)
			}
		} else {
			// join shard
			log.Printf("[INFO] shardmanager: decided to join subshard %d, nodeID: %d\n", shardID, nodeID)
			if err := it.JoinShard(ctx, shardID, nodeID); err != nil {
				return fmt.Errorf("JoinShard err: %w", err)
			}
			if err := it.StartCluster(map[uint64]string{}, true, it.fn, conf); err != nil {
				return fmt.Errorf("JoinSubShard %d start cluster err: %w", shardID, err)
			}
		}
	}
	return nil
}

func (it *ShardManager) JoinShard(ctx context.Context, shardID uint64, nodeID uint64) error {
	if _, err := it.client.NodeHost().AddNode(ctx, &service.ShardAddNodeRequest{
		ShardId: shardID,
		NodeId:  nodeID,
		Addr:    lo.ToPtr(it.cm.RaftAddress()),
	}); err != nil {
		return fmt.Errorf("NodeHost().AddNode(%d, %d) err: %w", shardID, nodeID, err)
	}
	return nil
}

func (it *ShardManager) WaitServeMasterShard(ctx context.Context) error {
	err := utils.RetryWithDelay(ctx, 20*60, time.Second*3, func() (bool, error) {
		if err := it.ServeMasterShard(ctx); err != nil {
			if errors.Is(err, ErrInitialMemberNotEnough) {
				log.Println("[INFO]", fmt.Sprintf("shardmanager: ErrInitialMemberNotEnough retrying"))
				return true, err
			}
			return false, err
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("WaitServeMasterShard timeout err: %w", err)
	}
	if err := it.WaitShardReady(ctx, cluster.MasterShardID, 60*30); err != nil {
		return fmt.Errorf("WaitShardReady(%d) err: %w", cluster.MasterShardID, err)
	}
	if err := it.ReconcileMigration(ctx, cluster.MasterShardID); err != nil {
		return fmt.Errorf("ReconcileMigration(%d) err: %w", cluster.MasterShardID, err)
	}

	return nil
}

func (it *ShardManager) StartCluster(initialMembers map[uint64]dragonboat.Target, join bool, create any, cfg config.Config) error {
	switch fn := create.(type) {
	case sm.CreateStateMachineFunc:
		return it.cm.NodeHost().StartCluster(initialMembers, join, fn, cfg)
	case sm.CreateConcurrentStateMachineFunc:
		return it.cm.NodeHost().StartConcurrentCluster(initialMembers, join, fn, cfg)
	case sm.CreateOnDiskStateMachineFunc:
		return it.cm.NodeHost().StartOnDiskCluster(initialMembers, join, fn, cfg)
	default:
		panic(fmt.Errorf("unknown StateMachineFunc type %T", fn))
	}
}

func (it *ShardManager) GetShardConfig(shardID uint64, nodeID uint64) config.Config {
	return config.Config{
		NodeID:              nodeID,
		ClusterID:           shardID,
		CheckQuorum:         true,
		ElectionRTT:         20,
		HeartbeatRTT:        2,
		SnapshotEntries:     100,
		CompactionOverhead:  0,
		OrderedConfigChange: false,
	}
}

func (it *ShardManager) RecoverShard(ctx context.Context, shardID uint64, nodeID uint64) error {
	if shardID == cluster.MasterShardID {
		return it.ServeMasterShard(ctx)
	} else {
		return it.ServeSubShard(ctx, shardID, nodeID, map[uint64]string{})
	}
}

func (it *ShardManager) RemoteQueryMigration(ctx context.Context, shardID uint64) (*service.MigrationState, error) {
	if resp, err := it.client.Spec().GetMigration(ctx, &service.QueryMigrationRequest{ShardId: shardID}); err != nil {
		return nil, fmt.Errorf("GetMigration(%d) err: %w", shardID, err)
	} else {
		return resp.State, nil
	}
}

func (it *ShardManager) PrepareMigration(ctx context.Context, shardID uint64) error {
	err := utils.RetryWithDelay(ctx, 10, time.Second*2, func() (bool, error) {
		if _, err := it.client.Spec().CreateMigration(ctx, &service.CreateMigrationRequest{ShardId: shardID}); err != nil {
			dcode := runtime.GetDragonboatErrorCode(service.GrpcErrorToDragonboatError(err))
			if dcode == runtime.ErrCodeMigrationExpired {
				return true, fmt.Errorf("CreateMigration(%d) err: %w", shardID, err)
			} else if dcode == runtime.ErrCodeAlreadyMigrating {
				log.Println("[WARN]", fmt.Errorf("PrepareMigration shard %d ErrCodeAlreadyMigrating err: %w", shardID, err))
				return false, nil
			}
			return false, fmt.Errorf("CreateMigration(%d) err: %w", shardID, err)
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("RetryWithDelay err: %w", err)
	}
	return nil
}

func (it *ShardManager) getMigrationNode(shardID uint64) *runtime.MigrationNode {
	nodeID := it.cm.GetNodeID(shardID)
	if nodeID == nil {
		panic(fmt.Errorf("call getMigrationNode before recovery shard %d", shardID))
	}
	nodes := lo.MapValues(it.cm.GetNodes(shardID), func(_ string, sid uint64) bool {
		return false
	})
	return &runtime.MigrationNode{
		NodeId:  *nodeID,
		Version: it.cm.Version().CodeVersion(),
		Nodes:   nodes,
	}
}

var ErrMigrationShardNotReady error = errors.New("ErrMigrationShardNotReady")

func (it *ShardManager) UpdateMigrationStatus(ctx context.Context, shardID uint64) error {
	node := it.getMigrationNode(shardID)
	if err := it.cm.Client().Migration(shardID).UpdateUpgradedNode(ctx, &runtime.DragonboatUpdateMigrationRequest{Node: node}, runtime.WithClientTimeout(time.Second*10)); err != nil {
		return fmt.Errorf("UpdateNodeMigration err: %w", err)
	}
	return nil
}

func (it *ShardManager) SaveMigration(ctx context.Context, shardID uint64) error {
	node := it.getMigrationNode(shardID)
	if err := it.cm.Client().Migration(shardID).UpdateSavedNode(ctx, &runtime.DragonboatCompleteMigrationRequest{Node: node}, runtime.WithClientTimeout(time.Second*10)); err != nil {
		return fmt.Errorf("UpdateSavedNode err: %w", err)
	}
	return nil
}

func (it *ShardManager) ReconcileMigration(ctx context.Context, shardID uint64) error {
	q, err := it.cm.Client().Migration(shardID).QueryMigration(ctx, &runtime.DragonboatQueryMigrationRequest{}, runtime.WithClientTimeout(time.Second*3))
	if err != nil {
		return fmt.Errorf("QueryMigration(%d) err: %w", shardID, err)
	}
	if q.View.Type == runtime.MigrationStateType_Upgrading {
		if err := it.UpdateMigrationStatus(ctx, shardID); err != nil {
			return fmt.Errorf("UpdateMigrationStatus(%d) err: %w", shardID, err)
		}
	} else if q.View.Type == runtime.MigrationStateType_Upgraded {
		if err := it.SaveMigration(ctx, shardID); err != nil {
			return fmt.Errorf("SaveMigration(%d) err: %w", shardID, err)
		}
	} else if q.View.Type == runtime.MigrationStateType_Empty || q.View.Type == runtime.MigrationStateType_Migrating || q.View.Type == runtime.MigrationStateType_Expired {
		if err := it.cm.Client().Migration(shardID).CheckVersion(ctx, it.cm.Version().CodeVersion(), runtime.WithClientTimeout(time.Second*3)); err != nil {
			return fmt.Errorf("CheckVersion(%d) err: %w", shardID, err)
		}
	}
	return nil
}

func (it *ShardManager) RecoveryShardWithMigration(ctx context.Context, shardID uint64, nodeID uint64, seconds int) error {
	mi, err := it.RemoteQueryMigration(ctx, shardID)
	if err == nil {
		if mi.Type == runtime.MigrationStateType_Empty.String() || mi.Type == runtime.MigrationStateType_Migrating.String() || mi.Type == runtime.MigrationStateType_Expired.String() {
			if mi.Version != it.cm.Version().CodeVersion() {
				if err := it.PrepareMigration(ctx, shardID); err != nil {
					return fmt.Errorf("PrepareMigration err: %w", err)
				}
			}
		}
	} else {
		log.Println("[INFO]", fmt.Sprintf("ShardManager: shard %d not health, recovery without migration", shardID))
	}
	if err := it.RecoverShard(ctx, shardID, nodeID); err != nil {
		return fmt.Errorf("RecoverShard err: %w", err)
	}
	if err := it.WaitShardReady(ctx, shardID, seconds); err != nil {
		return fmt.Errorf("WaitShardReady err: %w", err)
	}
	if err := it.ReconcileMigration(ctx, shardID); err != nil {
		if runtime.GetDragonboatErrorCode(err) == runtime.ErrCodeVersionNotMatch {
			replica := len(it.cm.GetNodes(shardID))
			if replica != 1 {
				panic(err)
			} else {
				log.Println("[WARN]", fmt.Errorf("ShardManager: shard %d skip migration code version cause of replica is %d", shardID, replica))
				if err := it.PrepareMigration(ctx, shardID); err != nil {
					return fmt.Errorf("PrepareMigration after recovery err: %w", err)
				}
			}
		} else {
			return fmt.Errorf("ReconcileMigration(%d) err: %w", shardID, err)
		}
	}

	return nil
}

func (it *ShardManager) ShardSpecChangingDispatcher(ctx context.Context, e *cluster.ShardSpecChangingEvent) error {
	switch e.Type {
	case cluster.ShardSpecChangingType_Adding:
		if err := it.ServeSubShard(ctx, e.ShardId, *e.CurrentNodeId, e.CurrentNodes); err != nil {
			return fmt.Errorf("ServeSubShard(%d) err: %w", e.ShardId, err)
		}
	case cluster.ShardSpecChangingType_Deleting:
		if err := it.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, false); err != nil {
			return fmt.Errorf("StopSubShard(%d, false) err: %w", e.ShardId, err)
		}
	case cluster.ShardSpecChangingType_Joining:
		if err := it.ServeSubShard(ctx, e.ShardId, *e.CurrentNodeId, nil); err != nil {
			return fmt.Errorf("ServeSubShard(%d) err: %w", e.ShardId, err)
		}
	case cluster.ShardSpecChangingType_Leaving:
		if err := it.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, true); err != nil {
			return fmt.Errorf("StopSubShard(%d, true) err: %w", e.ShardId, err)
		}
	// case cluster.ShardSpecChangingType_NodeJoining:
		// if it.cm.IsLeader(e.ShardId) {
			// nodeID := *e.CurrentNodeId
			// addr := e.CurrentNodes[nodeID]
			// if err := it.cm.Client().Raw(e.ShardId).AddNode(ctx, nodeID, addr, runtime.WithClientTimeout(time.Second*10)); err != nil {
				// return fmt.Errorf("Client.Raw().AddNode err: %w", err)
			// }
		// }
		// return nil
	// case cluster.ShardSpecChangingType_NodeLeaving:
		// if it.cm.IsLeader(e.ShardId) {
			// if err := it.cm.Client().Raw(e.ShardId).RemoveNode(ctx, *e.PreviousNodeId, runtime.WithClientTimeout(time.Second*10)); err != nil {
				// return fmt.Errorf("Client.Raw().RemoveNode err: %w", err)
			// }
		// }
		// return nil
	case cluster.ShardSpecChangingType_Cleanup:
		if err := it.StopSubShard(ctx, e.ShardId, *e.PreviousNodeId, false); err != nil {
			return fmt.Errorf("StopSubShard(%d) err: %w", e.ShardId, err)
		}
	case cluster.ShardSpecChangingType_Recovering:
		if err := it.RecoveryShardWithMigration(ctx, e.ShardId, *e.PreviousNodeId, 60); err != nil {
			return fmt.Errorf("RecoveryShardWithMigration(%d) err: %w", e.ShardId, err)
		}
	default:
		return fmt.Errorf("unknown event.Type: %T", e.Type)
	}
	return nil
}

func (it *ShardManager) IsShardReady(ctx context.Context, shardID uint64) error {
	if err := it.cm.Client().Void(shardID).VoidQuery(ctx, runtime.WithClientTimeout(time.Second)); err != nil {
		return fmt.Errorf("VoidQuery(%d) err: %w", shardID, err)
	}
	return nil
}

func (it *ShardManager) WaitShardReady(ctx context.Context, shardID uint64, seconds int) error {
	if err := utils.RetryWithDelay(ctx, seconds, time.Second, func() (bool, error) {
		if err := it.IsShardReady(ctx, shardID); err != nil {
			log.Println("[INFO]", fmt.Sprintf("Waiting shard %d ready", shardID))
			return true, fmt.Errorf("%w: %v", ErrMigrationShardNotReady, err)
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("RetryWithDelay err: %w", err)
	}
	return nil
}

func (it *ShardManager) LogEvents(events []*cluster.ShardSpecChangingEvent) {
	dict := lo.GroupBy(events, func(e *cluster.ShardSpecChangingEvent) cluster.ShardSpecChangingType {
		return e.Type
	})
	log.Println("[INFO]", fmt.Sprintf("ShardSpec changing events: Recovering: %d, Adding: %d, Cleanup: %d",
		len(dict[cluster.ShardSpecChangingType_Recovering]),
		len(dict[cluster.ShardSpecChangingType_Adding]),
		len(dict[cluster.ShardSpecChangingType_Cleanup]),
	))
}

func (it *ShardManager) RecoveryShards(ctx context.Context) error {
	if !it.cm.NodeHost().HasNodeInfo(cluster.MasterShardID, it.cm.MasterNodeID()) {
		log.Println("[INFO]", fmt.Sprintf("MasterShard not exists, skip RecoveryShards"))
		return nil
	}
	if err := it.RecoveryShardWithMigration(ctx, cluster.MasterShardID, it.cm.MasterNodeID(), 600); err != nil {
		return fmt.Errorf("RecoveryShardWithMigration MasterShard err: %w", err)
	}
	events, err := cluster.ReconcileSubShardChangingEvents(it.cm)
	if err != nil {
		return fmt.Errorf("ReconcileSubShardChangingEvents err: %w", err)
	}
	it.LogEvents(events)
	if _, err := utils.ParallelMap(events, func(e *cluster.ShardSpecChangingEvent) (uint64, error) {
		if err := it.ShardSpecChangingDispatcher(ctx, e); err != nil {
			return e.ShardId, fmt.Errorf("ShardSpecChangingDispatcher err: %w", err)
		}
		return e.ShardId, nil
	}); err != nil {
		log.Println("[WARN]", fmt.Errorf("ShardManager RecoverySubShards ParallelMap err: %w", err))
	}
	return nil
}

func (it *ShardManager) StopSubShard(ctx context.Context, shardID uint64, nodeID uint64, leave bool) error {
	if shardID == cluster.MasterShardID {
		return fmt.Errorf("stop master shard is prohibitted")
	}
	if leave {
		if err := it.cm.Client().Raw(shardID).RemoveNode(ctx, nodeID, runtime.WithClientTimeout(time.Second*10)); err != nil {
			if !errors.Is(err, dragonboat.ErrClusterNotFound) {
				return fmt.Errorf("Client.Raw.RemoveNode err: %w", err)
			}
		}
	}
	if err := it.cm.NodeHost().StopNode(shardID, nodeID); err != nil && !errors.Is(err, dragonboat.ErrClusterNotFound) {
		return fmt.Errorf("NodeHost().StopCluster(%d) err: %w", shardID, err)
	}
	if err := utils.RetryWithDelay(ctx, 3, time.Second, func() (bool, error) {
		if err := it.cm.NodeHost().RemoveData(shardID, nodeID); err != nil {
			if errors.Is(err, dragonboat.ErrClusterNotStopped) {
				return true, fmt.Errorf("NodeHost().RemoveData(%d, %d) err: %w", shardID, nodeID, err)
			}
			return false, fmt.Errorf("NodeHost().RemoveData(%d, %d) err: %w", shardID, nodeID, err)
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("RetryWithDelay max retry err: %w", err)
	}
	return nil
}

func (it *ShardManager) DrainLeader(ci dragonboat.ClusterInfo) error {
	if !ci.IsLeader {
		return nil
	}
	members := make(map[string]cluster.MemberNode, 0)
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		members[m.Meta.RaftAddress] = m
		return true
	})
	// get last startup member
	var targetNodeID uint64
	var maxStartupTimestamp int64
	for nodeID, addr := range ci.Nodes {
		if member, ok := members[addr]; ok {
			if member.Meta.RaftAddress == it.cm.RaftAddress() {
				continue
			}
			if member.Meta.StartupTimestamp > maxStartupTimestamp {
				targetNodeID = nodeID
				maxStartupTimestamp = member.Meta.StartupTimestamp
			}
		}
	}
	if maxStartupTimestamp == 0 {
		return fmt.Errorf("NoAvalibleNode")
	}
	if err := it.cm.NodeHost().RequestLeaderTransfer(ci.ClusterID, targetNodeID); err != nil {
		return fmt.Errorf("NodeHost().RequestLeaderTransfer(%d, %d) err: %w", ci.ClusterID, targetNodeID, err)
	}
	return nil
}

func (it *ShardManager) Drain() error {
	var errs error
	nhi := it.cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	for _, ci := range nhi.ClusterInfoList {
		if err := it.DrainLeader(ci); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
