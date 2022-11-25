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
	Snapshot() ISnapshotManager
	Version() IVersionManager

	MinClusterSize() int32
	NodeHostIndex() int32
	MasterNodeID() uint64
	ServiceName() string
	SnapshotBaseDir() string
	RecoverMode() bool

	ImportSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error
	Stop() error
}
