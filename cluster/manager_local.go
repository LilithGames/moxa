package cluster

import (
	"context"
	"fmt"
	"os"
	"time"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/goutils/vfs"

	"github.com/LilithGames/moxa/utils"
)

type LocalManager struct {
	config *Config
	nh     *utils.Provider[*dragonboat.NodeHost]
	ms     IMembers
	bus    *eventbus.EventBus
	cc     IClient
	ready  *Signal

	*NodeHostHelper
}

func NewLocalManager(ctx context.Context, config *Config) (Manager, error) {
	it := &LocalManager{
		config: config,
		nh:     utils.NewProvider[*dragonboat.NodeHost](nil),
		bus:    eventbus.NewEventBus(),
		ready:  NewSignal(false),
	}
	if err := it.startNodeHost(); err != nil {
		return nil, fmt.Errorf("startNodeHost err: %w", err)
	}
	it.NodeHostHelper = &NodeHostHelper{it.nh}
	hostName, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("os.Hostname err: %w", err)
	}
	meta := MemberMeta{
		HostName:         hostName,
		NodeHostId:       it.NodeHostID(),
		RaftAddress:      it.RaftAddress(),
		NodeHostIndex:    it.NodeHostIndex(),
		MasterNodeId:     it.MasterNodeID(),
		StartupTimestamp: time.Now().Unix(),
		Type:             MemberType_Dragonboat,
	}
	ms, err := NewMembers(meta, it.bus, config.MemberSeed)
	if err != nil {
		return nil, fmt.Errorf("NewMembers err: %w", err)
	}
	it.ms = ms
	it.cc = NewClient(it)

	if err := it.ms.UpdateNodeHostInfo(it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})); err != nil {
		return nil, fmt.Errorf("UpdateNodeHostInfo err : %w", err)
	}
	if err := utils.RetryWithDelay(ctx, 3, time.Second*3, func() (bool, error) {
		if err := it.ms.SyncState(); err != nil {
			return true, fmt.Errorf("SyncState err: %w", err)
		}
		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("RetryWithDelay err: %w", err)
	}
	return it, nil
}

func (it *LocalManager) startNodeHost() error {
	listener := &eventListener{bus: it.bus}
	nhconf := config.NodeHostConfig{
		DeploymentID:        it.config.DeploymentId,
		NodeHostDir:         it.config.NodeHostDir,
		AddressByNodeHostID: false,
		RTTMillisecond:      it.config.RttMillisecond,
		RaftAddress:         fmt.Sprintf("localhost:63000"),
		RaftEventListener:   listener,
		SystemEventListener: listener,
		EnableMetrics:       it.config.EnableMetrics,
	}
	if it.config.StorageType == StorageType_Memory {
		nhconf.Expert.FS = vfs.NewMem()
	}
	nh, err := dragonboat.NewNodeHost(nhconf)
	if err != nil {
		return fmt.Errorf("dragonboat.NewNodeHost err: %w", err)
	}
	it.nh.Set(nh)
	return nil
}

func (it *LocalManager) NodeHost() *dragonboat.NodeHost {
	return it.nh.Get()
}

func (it *LocalManager) Members() IMembers {
	return it.ms
}

func (it *LocalManager) Client() IClient {
	return it.cc
}

func (it *LocalManager) Config() *Config {
	return it.config
}

func (it *LocalManager) StartupReady() *Signal {
	return it.ready
}

func (it *LocalManager) EventBus() *eventbus.EventBus {
	return it.bus
}
func (it *LocalManager) MinClusterSize() int32 {
	return 1
}
func (it *LocalManager) MasterNodeID() uint64 {
	return it.NodeHostNum()
}

func (it *LocalManager) NodeHostIndex() int32 {
	return 0
}

func (it *LocalManager) ImportSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error {
	if err := it.LoadSnapshots(ctx, snapshots); err != nil {
		return fmt.Errorf("LoadSnapshots err: %w", err)
	}
	if err := it.startNodeHost(); err != nil {
		return fmt.Errorf("startNodeHost err: %w", err)
	}
	return nil
}

func (it *LocalManager) Stop() error {
	it.ms.Stop()
	it.nh.Get().Stop()
	return nil
}
