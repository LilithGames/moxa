package moxa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/service"
)

var ErrInitialMemberNotEnough error = errors.New("ErrInitialMemberNotEnough")

type ShardManager struct {
	cm     cluster.Manager
	client service.IClient
	fn     any
}

func NewShardManager(cm cluster.Manager, client service.IClient, fn interface{}) (*ShardManager, error) {
	return &ShardManager{cm, client, fn}, nil
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

func RetryWithDelay(ctx context.Context, max int, delay time.Duration, fn func() (bool, error)) error {
	var err error
	var retry bool
	for i := 0; i < max; i++ {
		retry, err = fn()
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

func ParallelMap[T comparable, R any](items []T, fn func(item T) (R, error)) (map[T]R, error) {
	ch := make(chan lo.Tuple2[T, R], len(items))
	ech := make(chan error, len(items))
	wg := sync.WaitGroup{}
	for _, item := range items {
		wg.Add(1)
		go func(arg T) {
			defer wg.Done()
			r, err := fn(arg)
			if err != nil {
				ech <- err
			} else {
				ch <- lo.T2(arg, r)
			}
		}(item)
	}
	wg.Wait()
	close(ch)
	close(ech)
	errs := lo.ChannelToSlice(ech)
	err := multierror.Append(nil, errs...).ErrorOrNil()
	r := lo.SliceToMap(lo.ChannelToSlice(ch), func(t lo.Tuple2[T, R]) (T, R) {
		return t.A, t.B
	})
	return r, err
}

func RemoteSnapshotSiri(src *service.RemoteSnapshot) *cluster.RemoteSnapshot {
	return &cluster.RemoteSnapshot{
		Index:      src.Index,
		Version:    src.Version,
		Path:       src.Path,
		NodeHostId: src.NodeHostId,
		Nodes:      src.Nodes,
	}
}

func (it *ShardManager) PrepareSnapshots(ctx context.Context) (map[uint64]*cluster.RemoteSnapshot, error) {
	nodeIDs := it.cm.GetLogNodeIDs()
	shardIDs := lo.MapToSlice(nodeIDs, func(shardID uint64, nodeID uint64) uint64 {
		return shardID
	})
	snapshots, err := ParallelMap(shardIDs, func(shardID uint64) (*cluster.RemoteSnapshot, error) {
		var snapshot *cluster.RemoteSnapshot
		if err := RetryWithDelay(ctx, 3, time.Second*5, func() (bool, error) {
			if resp, err := it.client.NodeHost().CreateSnapshot(ctx, &service.CreateSnapshotRequest{ShardId: shardID, Export: true, Leader: true}); err != nil {
				if e, ok := status.FromError(err); ok {
					if e.Code() == codes.AlreadyExists {
						if _, err := it.client.NodeHost().ShardNop(ctx, &service.ShardNopRequest{ShardId: shardID}); err != nil {
							return true, fmt.Errorf("ShardNop err: %w", err)
						}
						return true, fmt.Errorf("client.CreateSnapshot err: %w", err)
					} else if e.Code() == codes.Unavailable {
						return true, fmt.Errorf("client.CreateSnapshot err: %w", err)
					}
				}
				return false, fmt.Errorf("client.CreateSnapshot err: %w", err)
			} else {
				snapshot = RemoteSnapshotSiri(resp.Snapshot)
				return false, nil
			}
		}); err != nil {
			if it.cm.RecoverMode() {
				if _, herr := it.client.NodeHost().ShardHealthz(ctx, &service.ShardHealthzRequest{ShardId: shardID}); herr != nil {
					log.Println("[WARN]", fmt.Errorf("RecoverMode: skip prepare snapshot from unhealth shard %d, the may cause unconsistency in shard: %w, %v", shardID, err, herr))
					return nil, nil
				}
			}
			return nil, fmt.Errorf("Migration RetryWithDelay err: %w", err)
		}
		return snapshot, nil
	})
	return lo.OmitByValues(snapshots, []*cluster.RemoteSnapshot{nil}), err
}

func (it *ShardManager) ServeMasterShard(ctx context.Context) error {
	shardID := cluster.MasterShardID
	fn := sm.CreateStateMachineFunc(master_shard.NewStateMachine)
	conf := it.GetShardConfig(shardID, it.cm.MasterNodeID())
	if it.cm.NodeHost().HasNodeInfo(shardID, it.cm.MasterNodeID()) {
		// recovery
		log.Printf("[INFO] shardmanager: decided to recovery shard %d\n", shardID)
		if err := it.StartCluster(map[uint64]string{}, false, fn, conf); err != nil {
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
			if err := it.StartCluster(initials, false, fn, conf); err != nil {
				return fmt.Errorf("CreateShard %d err: %w", shardID, err)
			}
		} else {
			// join shard
			log.Printf("[INFO] shardmanager: decided to join shard %d\n", shardID)
			if _, err := it.client.NodeHost().AddNode(ctx, &service.ShardAddNodeRequest{ShardId: shardID, NodeIndex: it.cm.NodeHostIndex()}); err != nil {
				return fmt.Errorf("AddNode err: %w", err)
			}
			if err := it.StartCluster(map[uint64]string{}, true, fn, conf); err != nil {
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
			if _, err := it.client.NodeHost().AddNode(ctx, &service.ShardAddNodeRequest{
				ShardId:   shardID,
				NodeIndex: it.cm.NodeHostIndex(),
				NodeId:    &nodeID,
			}); err != nil {
				return fmt.Errorf("JoinSubShard %d add node err: %w", shardID, err)
			}
			if err := it.StartCluster(map[uint64]string{}, true, it.fn, conf); err != nil {
				return fmt.Errorf("JoinSubShard %d start cluster err: %w", shardID, err)
			}
		}
	}
	return nil
}

func (it *ShardManager) WaitServeMasterShard(ctx context.Context) error {
	err := RetryWithDelay(ctx, 20*60, time.Second*3, func() (bool, error) {
		if err := it.ServeMasterShard(ctx); err != nil {
			if errors.Is(err, ErrInitialMemberNotEnough) {
				return true, err
			}
			return false, err
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("WaitServeMasterShard timeout err: %w", err)
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
		ElectionRTT:         10,
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

func (it *ShardManager) Migrate(ctx context.Context) error {
	if it.cm.Version().Migrating() {
		log.Println("[INFO]", fmt.Sprintf("Migration start"))
		snapshots, err := it.PrepareSnapshots(ctx)
		if err != nil {
			return fmt.Errorf("PrepareSnapshots err: %w", err)
		}
		log.Println("[INFO]", fmt.Sprintf("PrepareSnapshots %d done", len(snapshots)))
		if err := it.cm.ImportSnapshots(ctx, snapshots); err != nil {
			return fmt.Errorf("ImportSnapshots err: %w", err)
		}
		log.Println("[INFO]", fmt.Sprintf("ImportSnapshots %d done", len(snapshots)))
	} else {
		log.Println("[INFO]", fmt.Sprintf("CodeVersion match with DataVersion, skip migration"))
	}
	return nil
}

func (it *ShardManager) RecoveryShards(ctx context.Context) error {
	nodeIDs := it.cm.GetLogNodeIDs()
	shardIDs := lo.MapToSlice(nodeIDs, func(shardID uint64, nodeID uint64) uint64 {
		return shardID
	})

	_, err := ParallelMap(shardIDs, func(shardID uint64) (uint64, error) {
		return shardID, it.RecoverShard(ctx, shardID, nodeIDs[shardID])
	})
	if err != nil {
		return err
	} else {
		it.cm.Version().Commit()
		return nil
	}
}

func (it *ShardManager) StopSubShard(ctx context.Context, shardID uint64, nodeID uint64, leave bool) error {
	if shardID == cluster.MasterShardID {
		return fmt.Errorf("stop master shard is prohibitted")
	}
	if leave {
		if _, err := it.client.NodeHost().RemoveNode(ctx, &service.ShardRemoveNodeRequest{ShardId: shardID, NodeIndex: it.cm.NodeHostIndex()}); err != nil {
			return fmt.Errorf("client.Shard().RemoveNode(%d) err: %w", shardID, err)
		}
	}
	if err := it.cm.NodeHost().StopNode(shardID, nodeID); err != nil && !errors.Is(err, dragonboat.ErrClusterNotFound) {
		return fmt.Errorf("NodeHost().StopCluster(%d) err: %w", shardID, err)
	}
	if err := RetryWithDelay(ctx, 3, time.Second, func() (bool, error) {
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
