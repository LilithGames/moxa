package master_shard

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func Test_IndexClient(t *testing.T) {
	core := NewCore()
	nodes := []*NodeCreateView{
		&NodeCreateView{NodeHostId: "nhid-1", Addr: "addr-1"},
		&NodeCreateView{NodeHostId: "nhid-2", Addr: "addr-2"},
		&NodeCreateView{NodeHostId: "nhid-3", Addr: "addr-3"},
	}
	resp, err := core.CreateShard(&CreateShardRequest{Name: "room-1", ProfileName: "profile-1", NodesView: &ShardNodesView{Nodes: nodes}, Labels: map[string]string{"room": "room-1"}})
	assert.Nil(t, err)
	_, err = core.CreateShard(&CreateShardRequest{Name: "room-2", ProfileName: "profile-1", NodesView: &ShardNodesView{Nodes: nodes}, Labels: map[string]string{"room": "room-2"}})
	assert.Nil(t, err)
	resp2, err := core.ListShards(&ListShardsRequest{IndexQuery: &ShardIndexQuery{Name: ShardSpecIndex_ShardID, Value: fmt.Sprintf("%d", resp.Shard.ShardId)}})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp2.Shards))
	resp3, err := core.ListShards(&ListShardsRequest{IndexQuery: &ShardIndexQuery{Name: ShardSpecIndex_LabelKey, Value: fmt.Sprintf("%s", "room")}})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp3.Shards))
	resp4, err := core.ListShards(&ListShardsRequest{IndexQuery: &ShardIndexQuery{Name: ShardSpecIndex_LabelKeyValue, Value: fmt.Sprintf("%s", "room:room-2")}})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp4.Shards))
}
