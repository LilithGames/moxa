package cluster

import (
	"context"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3"
)

const MasterShardID uint64 = 0

type Manager interface {
	INodeHostHelper

	NodeHost() *dragonboat.NodeHost
	Members() IMembers
	EventBus() *eventbus.EventBus
	Client() IClient
	Version() IVersionManager
	Config() *Config
	StartupReady() *Signal

	MinClusterSize() int32
	NodeHostIndex() int32
	MasterNodeID() uint64
	ServiceName() string

	ImportSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error
	Stop() error
}
