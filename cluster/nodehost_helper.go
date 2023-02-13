package cluster

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/tools"
)

type INodeHostHelper interface {
	NodeHostID() string
	NodeHostNum() uint64
	GetNodeIDs() map[uint64]uint64
	GetNodeID(shardID uint64) *uint64
	IsLeader(shardID uint64) bool
	GetLeaderID(shardID uint64) *uint64
	RaftAddress() string
	GetNodes(shardID uint64) map[uint64]string
	GetLogNodeIDs() map[uint64]uint64
	GetLogNodeID(shardID uint64) *uint64
}

type NodeHostHelper struct {
	nh *Provider[*dragonboat.NodeHost]
}

func (it *NodeHostHelper) GetNodeIDs() map[uint64]uint64 {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	shards := make(map[uint64]uint64, len(nhi.ClusterInfoList))
	for _, s := range nhi.ClusterInfoList {
		shards[s.ClusterID] = s.NodeID
	}
	return shards
}

func (it *NodeHostHelper) NodeHostID() string {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	return nhi.NodeHostID
}

func (it *NodeHostHelper) GetNodeID(shardID uint64) *uint64 {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	for _, ci := range nhi.ClusterInfoList {
		if ci.ClusterID == shardID {
			return &ci.NodeID
		}
	}
	return nil

}

func (it *NodeHostHelper) GetLogNodeIDs() map[uint64]uint64 {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})
	nodeIDs := make(map[uint64]uint64, len(nhi.LogInfo))
	for _, ni := range nhi.LogInfo {
		nodeIDs[ni.ClusterID] = ni.NodeID
	}
	return nodeIDs
}

func (it *NodeHostHelper) GetLogNodeID(shardID uint64) *uint64 {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})
	for _, ni := range nhi.LogInfo {
		if ni.ClusterID == shardID {
			return &ni.NodeID
		}
	}
	return nil
}

func (it *NodeHostHelper) GetLeaderID(shardID uint64) *uint64 {
	leaderNodeID, valid, err := it.nh.Get().GetLeaderID(shardID)
	if err != nil || !valid {
		return nil
	}
	return &leaderNodeID
}

func (it *NodeHostHelper) IsLeader(shardID uint64) bool {
	leaderNodeID, valid, err := it.nh.Get().GetLeaderID(shardID)
	if err != nil || !valid {
		return false
	}
	nodeID := it.GetNodeID(shardID)
	if nodeID == nil {
		return false
	}
	return leaderNodeID == *nodeID
}

func (it *NodeHostHelper) NodeHostNum() uint64 {
	return parseNodeHostID(it.NodeHostID())
}

func (it *NodeHostHelper) RaftAddress() string {
	return it.nh.Get().NodeHostConfig().RaftAddress
}

func (it *NodeHostHelper) GetNodes(shardID uint64) map[uint64]string {
	nhi := it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	for _, ci := range nhi.ClusterInfoList {
		if ci.ClusterID == shardID {
			return ci.Nodes
		}
	}
	return nil
}

func (it *NodeHostHelper) LoadSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	nh := it.nh.Get()
	conf := nh.NodeHostConfig()
	nodeIDs := it.GetLogNodeIDs()
	nh.Stop()
	for shardID, snapshot := range snapshots {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		nodeID, ok := nodeIDs[shardID]
		if !ok {
			return fmt.Errorf("snapshot shard %d info not found", shardID)
		}
		if err := tools.ImportSnapshot(conf, snapshot.Path, snapshot.Nodes, nodeID); err != nil {
			return fmt.Errorf("tools.ImportSnapshot(%s) err: %w", snapshot, err)
		}
		log.Println("[INFO]", fmt.Sprintf("Import snapshot shard: %d nodeID: %d snapshot: %v success", shardID, nodeID, snapshot))
	}
	return nil
}

func getNodeHostIDString(nhid uint64) string {
	return fmt.Sprintf("nhid-%d", nhid)
}

func parseNodeHostID(nhid string) uint64 {
	items := strings.Split(nhid, "-")
	if len(items) != 2 || items[0] != "nhid" {
		panic(fmt.Errorf("invalid nodeHostID %s", nhid))
	}
	id, err := strconv.ParseUint(items[len(items)-1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid nodeHostID %s, ParseUint err: %w", nhid, err))
	}
	return id
}
