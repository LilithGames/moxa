package master_shard

import (
	"fmt"
)


type ShardSpecIndexSource struct  {
	server IMasterShardDragonboatServer
}

func NewShardSpecIndexSource(server IMasterShardDragonboatServer) IIndexSource[string, *ShardSpec] {
	return &ShardSpecIndexSource{server: server}
}

func (it *ShardSpecIndexSource) List() []*ShardSpec {
	resp, err := it.server.ListShards(&ListShardsRequest{})
	if err != nil {
		panic(fmt.Errorf("ListShards should not err without any query on server side"))
	}
	return resp.Shards
}

func (it *ShardSpecIndexSource) Get(key string) (*ShardSpec, error) {
	resp, err := it.server.GetShard(&GetShardRequest{Name: key})
	if err != nil {
		return nil, fmt.Errorf("server.GetShard(%s) err: %w", key, err)
	}
	return resp.Shard, nil
}
func (it *ShardSpecIndexSource) Key(item *ShardSpec) string {
	return item.ShardName
}
