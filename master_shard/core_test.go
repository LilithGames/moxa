package master_shard

import (
	"testing"
	"fmt"

	assert "github.com/stretchr/testify/require"
)

func TestCore(t *testing.T) {
	core := NewCore()
	nodes := []*NodeCreateView{
		&NodeCreateView{NodeHostId: "nhid-1", Addr: "addr-1"},
		&NodeCreateView{NodeHostId: "nhid-2", Addr: "addr-2"},
		&NodeCreateView{NodeHostId: "nhid-3", Addr: "addr-3"},
	}
	resp, err := core.CreateShard(&CreateShardRequest{Name: "auto-1", Nodes: nodes})
	assert.Nil(t, err)
	assert.Equal(t, uint64(4), resp.Shard.LastNodeId)
	var replica2 int32 = 2
	var replica3 int32 = 3
	resp2, err := core.UpdateShard(&UpdateShardRequest{Name: "auto-1", Nodes: nodes, Replica: &replica2})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", Nodes: nodes, Replica: &replica3})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", Nodes: nodes, Replica: &replica2})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp2, err = core.UpdateShard(&UpdateShardRequest{Name: "auto-1", Nodes: nodes, Replica: &replica3})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), resp2.Updated)
	resp3, err := core.GetShard(&GetShardRequest{Name: "auto-1"})
	assert.Nil(t, err)
	assert.Equal(t, uint64(6), resp3.Shard.LastNodeId)
	fmt.Printf("%+v\n", resp3.Shard)
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
