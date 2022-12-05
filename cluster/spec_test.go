package cluster

import (
	"testing"
	"fmt"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/raftio"
	"github.com/stretchr/testify/assert"

	"github.com/LilithGames/moxa/master_shard"
)

func TestDiff_Add(t *testing.T) {
	addr := "addr-3"
	spec := []*master_shard.ShardSpec{
		&master_shard.ShardSpec{
			ShardName: "auto-1",
			ShardId: 10000000,
			Replica: 3,
			Initials: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
				"nhid-2": &master_shard.Node{NodeHostId: "nhid-2", NodeId: 2, Addr: "addr-2"},
				"nhid-3": &master_shard.Node{NodeHostId: "nhid-3", NodeId: 3, Addr: "addr-3"},
			},
			Nodes: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
				"nhid-2": &master_shard.Node{NodeHostId: "nhid-2", NodeId: 2, Addr: "addr-2"},
				"nhid-3": &master_shard.Node{NodeHostId: "nhid-3", NodeId: 3, Addr: "addr-3"},
			},
			LastNodeId: 4,
		},
	}
	status := []dragonboat.ClusterInfo{}
	persist := []raftio.NodeInfo{}
	events := DiffSubShards(addr, spec, status, persist)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, ShardSpecChangingType_Adding, events[0].Type)
	
	fmt.Printf("%+v\n", events)
}

func TestDiff_Join(t *testing.T) {
	addr := "addr-4"
	spec := []*master_shard.ShardSpec{
		&master_shard.ShardSpec{
			ShardName: "auto-1",
			ShardId: 10000000,
			Replica: 3,
			Initials: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
				"nhid-2": &master_shard.Node{NodeHostId: "nhid-2", NodeId: 2, Addr: "addr-2"},
				"nhid-3": &master_shard.Node{NodeHostId: "nhid-3", NodeId: 3, Addr: "addr-3"},
			},
			Nodes: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
				"nhid-2": &master_shard.Node{NodeHostId: "nhid-2", NodeId: 2, Addr: "addr-2"},
				"nhid-4": &master_shard.Node{NodeHostId: "nhid-4", NodeId: 4, Addr: "addr-4"},
			},
			LastNodeId: 5,
		},
	}
	status := []dragonboat.ClusterInfo{}
	persist := []raftio.NodeInfo{}
	events := DiffSubShards(addr, spec, status, persist)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, ShardSpecChangingType_Joining, events[0].Type)
	
	fmt.Printf("%+v\n", events)
}

func TestDiff_NodeJoin(t *testing.T) {
	addr := "addr-1"
	spec := []*master_shard.ShardSpec{
		&master_shard.ShardSpec{
			ShardName: "auto-1",
			ShardId: 10000000,
			Replica: 3,
			Nodes: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
				"nhid-2": &master_shard.Node{NodeHostId: "nhid-2", NodeId: 2, Addr: "addr-2"},
				"nhid-3": &master_shard.Node{NodeHostId: "nhid-3", NodeId: 3, Addr: "addr-3"},
			},
			LastNodeId: 4,
		},
	}
	status := []dragonboat.ClusterInfo{
		dragonboat.ClusterInfo{
			ClusterID: 10000000,
			NodeID: 1, 
			Nodes: map[uint64]string{
				1: "addr-1",
			},
			IsLeader: true,
		},
	}
	persist := []raftio.NodeInfo{
		raftio.NodeInfo{ClusterID: 10000000, NodeID: 1},
	}
	events := DiffSubShards(addr, spec, status, persist)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, ShardSpecChangingType_NodeJoining, events[0].Type)
	assert.Equal(t, ShardSpecChangingType_NodeJoining, events[1].Type)
	
	fmt.Printf("%+v\n", events)
}

func TestDiff_NodeLeave(t *testing.T) {
	addr := "addr-1"
	spec := []*master_shard.ShardSpec{
		&master_shard.ShardSpec{
			ShardName: "auto-1",
			ShardId: 10000000,
			Replica: 1,
			Nodes: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
			},
			LastNodeId: 4,
		},
	}
	status := []dragonboat.ClusterInfo{
		dragonboat.ClusterInfo{
			ClusterID: 10000000,
			NodeID: 1, 
			Nodes: map[uint64]string{
				1: "addr-1",
				2: "addr-2",
			},
			IsLeader: true,
		},
	}
	persist := []raftio.NodeInfo{
		raftio.NodeInfo{ClusterID: 10000000, NodeID: 1},
	}
	events := DiffSubShards(addr, spec, status, persist)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, ShardSpecChangingType_NodeLeaving, events[0].Type)
	
	fmt.Printf("%+v\n", events)
}

func TestDiff_NodeLeaveCleanup(t *testing.T) {
	addr := "addr-2"
	spec := []*master_shard.ShardSpec{
		&master_shard.ShardSpec{
			ShardName: "auto-1",
			ShardId: 10000000,
			Replica: 1,
			Nodes: map[string]*master_shard.Node{
				"nhid-1": &master_shard.Node{NodeHostId: "nhid-1", NodeId: 1, Addr: "addr-1"},
			},
			LastNodeId: 4,
		},
	}
	status := []dragonboat.ClusterInfo{
	}
	persist := []raftio.NodeInfo{
		raftio.NodeInfo{ClusterID: 10000000, NodeID: 2},
	}
	events := DiffSubShards(addr, spec, status, persist)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, ShardSpecChangingType_Cleanup, events[0].Type)
	
	fmt.Printf("%+v\n", events)

}
