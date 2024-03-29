package cluster

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/goutils/vfs"
	"github.com/miekg/dns"

	"github.com/LilithGames/moxa/utils"
)

type KubernetesManager struct {
	config *Config
	nh     *utils.Provider[*dragonboat.NodeHost]
	ms     IMembers
	bus    *eventbus.EventBus
	cc     IClient
	ready  *Signal

	*NodeHostHelper
}

func NewKubernetesManager(ctx context.Context, config *Config) (Manager, error) {
	if err := validateKubernetesEnvs(); err != nil {
		return nil, fmt.Errorf("validateKubernetesEnvs err: %w", err)
	}
	it := &KubernetesManager{
		config: config,
		nh:     utils.NewProvider[*dragonboat.NodeHost](nil),
		bus:    eventbus.NewEventBus(),
		ready:  NewSignal(false),
	}
	if err := it.ensureRaftAddressResolved(ctx); err != nil {
		return nil, fmt.Errorf("ensureRaftAddressResolved err: %w", err)
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

func (it *KubernetesManager) startNodeHost() error {
	listener := &eventListener{bus: it.bus}
	nhconf := config.NodeHostConfig{
		DeploymentID:        it.config.DeploymentId,
		NodeHostDir:         it.config.NodeHostDir,
		AddressByNodeHostID: false,
		RTTMillisecond:      it.config.RttMillisecond,
		RaftAddress:         fmt.Sprintf("%s.%s.%s.svc.cluster.local:63000", os.Getenv("POD_NAME"), os.Getenv("POD_SERVICENAME"), os.Getenv("POD_NAMESPACE")),
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

func (it *KubernetesManager) ensureRaftAddressResolved(ctx context.Context) error {
	domain := fmt.Sprintf("%s.%s.%s.svc.cluster.local.", os.Getenv("POD_NAME"), os.Getenv("POD_SERVICENAME"), os.Getenv("POD_NAMESPACE"))
	podIP := os.Getenv("POD_IP")
	ip := net.ParseIP(podIP)
	if ip == nil {
		return fmt.Errorf("invalid POD_IP: %s", podIP)
	}
	config, err := GetDefaultDnsClientConfig()
	if err != nil {
		return fmt.Errorf("GetDefaultDnsClientConfig err: %w", err)
	}
	client := &dns.Client{Timeout: time.Second}
	timer := time.NewTimer(time.Second)
	for i := 0; i < 10; {
		ips, err := DnsResolveRecordA(client, config, domain)
		if err != nil {
			log.Println("[WARN]", fmt.Errorf("KubernetesManager DnsResolveRecordA err: %w", err))
			i = 0
		} else if len(ips) == 0 {
			log.Println("[WARN]", fmt.Errorf("KubernetesManager DnsResolveRecordA empty result"))
			i = 0
		} else if !ips[0].Equal(ip) {
			log.Println("[WARN]", fmt.Errorf("KubernetesManager resolve expect: %s, actual: %s", ip.String(), ips[0].String()))
			i = 0
		} else {
			i++
		}
		utils.ResetTimer(timer, time.Second)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return fmt.Errorf("waiting raft address resolve ctx err: %w", ctx.Err())
		}
	}
	return nil
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

func (it *KubernetesManager) Config() *Config {
	return it.config
}

func (it *KubernetesManager) StartupReady() *Signal {
	return it.ready
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

func (it *KubernetesManager) NodeHostIndex() int32 {
	return getPodNodeHostIndex(os.Getenv("POD_NAME"))
}

func (it *KubernetesManager) ImportSnapshots(ctx context.Context, snapshots map[uint64]*RemoteSnapshot) error {
	if err := it.LoadSnapshots(ctx, snapshots); err != nil {
		return fmt.Errorf("LoadSnapshots err: %w", err)
	}
	if err := it.startNodeHost(); err != nil {
		return fmt.Errorf("startNodeHost err: %w", err)
	}
	return nil
}

func (it *KubernetesManager) Stop() error {
	it.ms.Stop()
	it.nh.Get().Stop()
	time.Sleep(time.Second * 3)
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
