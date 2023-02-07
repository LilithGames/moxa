package master_shard

import (
	"testing"
	"fmt"

	"github.com/samber/lo"
	assert "github.com/stretchr/testify/require"
)

func TestCore(t *testing.T) {
	core := NewCore()
	nodes := []*NodeCreateView{
		{NodeHostId: "nhid-1", Addr: "addr-1"},
		{NodeHostId: "nhid-2", Addr: "addr-2"},
		{NodeHostId: "nhid-3", Addr: "addr-3"},
	}
	resp, err := core.CreateShard(&CreateShardRequest{Name: "auto-1", ProfileName: "profile-1", NodesView: &ShardNodesView{Nodes: nodes}})
	assert.Nil(t, err)
	assert.Equal(t, uint64(4), resp.Shard.LastNodeId)

	query := "'nhid-1' in it.Nodes"
	resp4, err := core.ListShards(&ListShardsRequest{Query: &query})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp4.Shards))

	var replica2 int32 = 2
	var replica3 int32 = 3
	resp2, err := core.UpdateShard(&UpdateShardRequest{Name: "auto-1", NodesView: &ShardNodesView{Nodes: nodes, Replica: &replica2}})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", NodesView: &ShardNodesView{Nodes: nodes, Replica: &replica3}})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", NodesView: &ShardNodesView{Nodes: nodes, Replica: &replica2}})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", NodesView: &ShardNodesView{Nodes: nodes, Replica: &replica3}})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)

	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", UpdateView: &ShardUpdateView{Labels: map[string]string{"room": "room-1"}}})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)

	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", UpdateView: &ShardUpdateView{StateVersion: lo.ToPtr(uint64(0)), Labels: map[string]string{"room": "room-1"}}})
	assert.NotNil(t, err)

	resp3, err := core.GetShard(&GetShardRequest{Name: "auto-1"})
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), resp3.Shard.LastNodeId)

	resp, err = core.CreateShard(&CreateShardRequest{Name: "auto-2", ProfileName: "profile-1", NodesView: &ShardNodesView{Nodes: nodes}})
	assert.Nil(t, err)

	resp5, err := core.ListShards(&ListShardsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp5.Shards))

	fmt.Printf("%+v\n", resp5.Shards)
}

func TestNodeSpec(t *testing.T) {
	core := NewCore()
	node := &NodeSpec{NodeHostId: "nhid-1", Labels: map[string]string{"k1": "v1"}}
	resp, err := core.UpdateNode(&UpdateNodeRequest{Node: node})
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), resp.Updated)
	resp2, err := core.GetNode(&GetNodeRequest{NodeHostId: "nhid-1"})
	assert.Nil(t, err)
	assert.Equal(t, "nhid-1", resp2.Node.NodeHostId)
	resp3, err := core.GetNode(&GetNodeRequest{NodeHostId: "nhid-2"})
	assert.NotNil(t, err)
	_ = resp3

	node1 := resp2.Node
	resp6, err := core.UpdateNode(&UpdateNodeRequest{Node: node1})
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), resp6.Updated)

	node1.Labels["k2"] = "v2"
	resp4, err := core.UpdateNode(&UpdateNodeRequest{Node: node1})
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), resp4.Updated)

	node.Labels["k2"] = "v3"
	resp5, err := core.UpdateNode(&UpdateNodeRequest{Node: node1})
	assert.NotNil(t, err)
	_ = resp5
}
