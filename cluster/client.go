package cluster

import (
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/LilithGames/moxa/master_shard"
)

type IClient interface {
	MasterShard() master_shard.IMasterShardDragonboatClient
	Migration(shardID uint64) runtime.IMigrationDragonboatClient
	Void(shardID uint64) runtime.IVoidDragonboatClient
	Raw(shardID uint64) runtime.IDragonboatClient
}

type Client struct {
	cm Manager
}

func NewClient(cm Manager) IClient {
	return &Client{cm}
}

func (it *Client) MasterShard() master_shard.IMasterShardDragonboatClient {
	client := runtime.NewDragonboatClient(it.cm.NodeHost(), MasterShardID)
	return master_shard.NewMasterShardDragonboatClient(client)
}

func (it *Client) Migration(shardID uint64) runtime.IMigrationDragonboatClient {
	return runtime.NewMigrationDragonboatClient(it.Raw(shardID))
}

func (it *Client) Void(shardID uint64) runtime.IVoidDragonboatClient {
	return runtime.NewVoidDragonboatClient(it.Raw(shardID))
}

func (it *Client) Raw(shardID uint64) runtime.IDragonboatClient {
	return runtime.NewDragonboatClient(it.cm.NodeHost(), shardID)
}
