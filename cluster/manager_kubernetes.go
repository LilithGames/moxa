package cluster

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"log"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/tools"
	"github.com/lni/goutils/vfs"
)

type KubernetesManager struct {
	nh      *Provider[*dragonboat.NodeHost]
	ms      IMembers
	bus     *eventbus.EventBus
	cc      IClient
	store   IStorage
	version IVersionManager

	*NodeHostHelper
}

func NewKubernetesManager() (Manager, error) {
	if err := validateKubernetesEnvs(); err != nil {
		return nil, fmt.Errorf("validateKubernetesEnvs err: %w", err)
	}
	store := NewLocalStorage("/data/storage")
	version, err := NewVersionManager(store)
	if err != nil {
		return nil, fmt.Errorf("NewVersionManager err: %w", err)
	}
	log.Println("[INFO]", fmt.Sprintf("DataVersion: %s", version.DataVersion()))
	log.Println("[INFO]", fmt.Sprintf("CodeVersion: %s", version.CodeVersion()))
	it := &KubernetesManager{
		nh:      NewProvider[*dragonboat.NodeHost](nil),
		bus:     eventbus.NewEventBus(),
		store:   store,
		version: version,
	}
	if err := it.startNodeHost(); err != nil {
		return nil, fmt.Errorf("startNodeHost err: %w", err)
	}
	it.NodeHostHelper = &NodeHostHelper{it.nh}
	meta := MemberMeta{
		HostName:         os.Getenv("POD_NAME"),
		NodeHostId:       it.NodeHostID(),
		RaftAddress:      it.RaftAddress(),
		NodeHostIndex:    it.NodeHostIndex(),
		MasterNodeId:     it.MasterNodeID(),
		StartupTimestamp: time.Now().Unix(),
		CodeHash:         it.version.CodeVersion(),
		Type:             MemberType_Dragonboat,
	}
	ms, err := NewMembers(meta)
	if err != nil {
		return nil, fmt.Errorf("NewMembers err: %w", err)
	}
	it.ms = ms
	it.cc = NewClient(it)

	if err := it.ms.UpdateNodeHostInfo(it.nh.Get().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})); err != nil {
		return nil, fmt.Errorf("UpdateNodeHostInfo err : %w", err)
	}
	if err := it.ms.SyncState(); err != nil {
		return nil, fmt.Errorf("SyncState err: %w", err)
	}
	return it, nil
}

func (it *KubernetesManager) startNodeHost() error {
	listener := &eventListener{bus: it.bus}
	nhconf := config.NodeHostConfig{
		NodeHostDir:         "/data/nodehost",
		AddressByNodeHostID: false,
		RTTMillisecond:      100,
		RaftAddress:         fmt.Sprintf("%s.%s.%s.svc.cluster.local:63000", os.Getenv("POD_NAME"), os.Getenv("POD_SERVICENAME"), os.Getenv("POD_NAMESPACE")),
		RaftEventListener:   listener,
		SystemEventListener: listener,
	}
	nh, err := dragonboat.NewNodeHost(nhconf)
	if err != nil {
		return fmt.Errorf("dragonboat.NewNodeHost err: %w", err)
	}
	it.nh.Set(nh)
	return nil
}

func (it *KubernetesManager) Version() IVersionManager {
	return it.version
}

func (it *KubernetesManager) Snapshot() ISnapshotManager {
	return NewSnapshotManager(it.SnapshotBaseDir(), it.NodeHostID(), vfs.Default)
}
func (it *KubernetesManager) NodeHost() *dragonboat.NodeHost {
	return it.nh.Get()
}

func (it *KubernetesManager) Members() IMembers {
	return it.ms
}

func (it *KubernetesManager) Client() IClient {
	return it.cc
}

func (it *KubernetesManager) RecoverMode() bool {
	return strings.ToLower(os.Getenv("RECOVER_MODE")) == "true"
}

func (it *KubernetesManager) EventBus() *eventbus.EventBus {
	return it.bus
}
func (it *KubernetesManager) MinClusterSize() int32 {
	return 3
}
func (it *KubernetesManager) MasterNodeID() uint64 {
	return it.NodeHostNum()
}
func (it *KubernetesManager) ServiceName() string {
	return os.Getenv("POD_SERVICENAME")
}
func (it *KubernetesManager) SnapshotBaseDir() string {
	return fmt.Sprintf("%s/snapshot", os.Getenv("POD_SHAREDIR"))
}

func (it *KubernetesManager) NodeHostIndex() int32 {
	return getPodNodeHostIndex(os.Getenv("POD_NAME"))
}

// TODO: concurrency import
func (it *KubernetesManager) ImportSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error {
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
	}
	if err := it.startNodeHost(); err != nil {
		return fmt.Errorf("startNodeHost err: %w", err)
	}
	return nil
}

func (it *KubernetesManager) Stop() error {
	it.ms.Stop()
	it.nh.Get().Stop()
	return nil
}

func validateKubernetesEnvs() error {
	if os.Getenv("POD_IP") == "" {
		return fmt.Errorf("needs env POD_IP")
	}
	if os.Getenv("POD_NAMESPACE") == "" {
		return fmt.Errorf("needs env POD_NAMESPACE")
	}
	if os.Getenv("POD_NAME") == "" {
		return fmt.Errorf("needs env POD_NAME")
	}
	if os.Getenv("POD_SERVICENAME") == "" {
		return fmt.Errorf("needs env POD_SERVICENAME")
	}
	if os.Getenv("POD_SHAREDIR") == "" {
		return fmt.Errorf("needs env POD_SHAREDIR")
	}

	return nil
}

func getPodNodeHostIndex(podname string) int32 {
	items := strings.Split(podname, "-")
	if len(items) < 2 {
		panic("invalid podname: " + podname)
	}
	index, err := strconv.ParseInt(items[len(items)-1], 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(index)
}
