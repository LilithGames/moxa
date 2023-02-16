package cluster

import (
	"fmt"
	"time"
	"log"

	"github.com/hashicorp/memberlist"
	"google.golang.org/protobuf/proto"
	"github.com/lni/dragonboat/v3"
	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/goutils/syncutil"

	"github.com/LilithGames/moxa/utils"
)

type MemberNode struct {
	Node *memberlist.Node
	Meta *MemberMeta
	State *MemberState
}


type IMembers interface {
	IMemberStateManager
	UpdateNodeHostInfo(nhi *dragonboat.NodeHostInfo) error
	SyncState() error
	Foreach(func(MemberNode) bool)
	ChangesSince(version uint64) (chan utils.Stream[utils.SyncStateView[*MemberState]], chan struct{})
	Nums() int
	Notify(notify *MemberNotify, opts NotifyOptions) error
	Stop() error
}

type Members struct {
	seed []string

	meta *MemberMeta
	ml   *memberlist.Memberlist
	*MemberStateManager
	bus *eventbus.EventBus
	stopper *syncutil.Stopper
}

func NewMembers(meta MemberMeta, bus *eventbus.EventBus, seed []string) (IMembers, error) {
	it := &Members{
		seed: seed,
		meta: &meta,
		bus:  bus,
		stopper: syncutil.NewStopper(),
		MemberStateManager: NewMemberStateManager(meta.NodeHostId, bus),
	}
	conf := memberlist.DefaultLANConfig()
	conf.Name = meta.NodeHostId
	conf.Delegate = it
	conf.Logger = log.Default()
	conf.PushPullInterval = time.Second * 30
	ml, err := memberlist.Create(conf)
	if err != nil {
		return nil, fmt.Errorf("memberlist.Create err %w", err)
	}
	it.ml = ml
	return it, nil
}

func (it *Members) SyncState() error {
	if _, err := it.ml.Join(it.seed); err != nil {
		return fmt.Errorf("memberlist.Join err %w", err)
	}
	return nil
}

func (it *Members) UpdateNodeHostInfo(nhi *dragonboat.NodeHostInfo) error {
	if nhi == nil {
		return fmt.Errorf("nhi is required")
	}
	shards := make(map[uint64]*MemberShard)
	logShards := make(map[uint64]uint64)
	for _, ci := range nhi.ClusterInfoList {
		shards[ci.ClusterID] = &MemberShard {
			ShardId: ci.ClusterID,
			NodeId: ci.NodeID,
			Nodes: ci.Nodes,
			IsLeader: ci.IsLeader,
			IsObserver: ci.IsObserver,
			IsWitness: ci.IsWitness,
			Pending: ci.Pending,
		}
	}
	for _, ni := range nhi.LogInfo {
		logShards[ni.ClusterID] = ni.NodeID
	}
	state := MemberState{
		NodeHostId: nhi.NodeHostID,
		Meta: it.meta,
		Shards: shards,
		LogShards: logShards,
	}
	if err := it.SetMemberState(state); err != nil {
		return fmt.Errorf("SetMemberState err: %w", err)
	}
	return nil
}

func (it *Members) Foreach(fn func(MemberNode) bool) {
	mstate, _ := it.GetMemberStateList()
	for _, member := range it.ml.Members() {
		meta := &MemberMeta{}
		if err := proto.Unmarshal(member.Meta, meta); err != nil {
			panic(fmt.Errorf("proto.Unmarshal member meta err: %w", err))
		}
		node := MemberNode{Node: member, Meta: meta}
		if state, ok := mstate[meta.NodeHostId]; ok {
			node.State = state
		}
		if !fn(node) {
			return
		}
	}
}

func (it *Members) ChangesSince(version uint64) (chan utils.Stream[utils.SyncStateView[*MemberState]], chan struct{}) {
	stream := make(chan utils.Stream[utils.SyncStateView[*MemberState]], 0)
	done := make(chan struct{}, 0)
	w := utils.NewStreamWriter(stream)
	it.stopper.RunWorker(func() {
		defer w.Done()
		var curr uint64
		changed, sub := it.bus.Subscribe(EventTopic_MemberStateChanged.String())
		defer sub.Close()
		items, err := it.MemberStateManager.Since(version)
		if err != nil {
			w.PutError(fmt.Errorf("MemberStateManager.Since(%d), err: %w", version, err))
			return
		}
		for _, item := range items {
			w.PutItem(item)
			if item.Version > curr {
				curr = item.Version
			}
		}
		for {
			select {
			case ev := <-changed:
				e := ev.Data.(*MemberStateChangedEvent)
				if e.Version > curr {
					items, err := it.MemberStateManager.Get(e.Version)
					if err != nil {
						w.PutError(fmt.Errorf("MemberStateManager.Get(%d) err: %w", e.Version, err))
						return
					}
					for _, item := range items {
						w.PutItem(item)
					}
					curr = e.Version
				}
			case <-done:
				return 
			case <-it.stopper.ShouldStop():
				w.PutError(fmt.Errorf("Members closed"))
				return
			}
		}
	})
	return stream, done
}

func (it *Members) Nums() int {
	return it.ml.NumMembers()
}

func (it *Members) Stop() error {
	it.stopper.Stop()
	it.ml.Leave(time.Second)
	return it.ml.Shutdown()
}

func (it *Members) NodeMeta(limit int) []byte {
	bs, err := proto.Marshal(it.meta)
	if err != nil {
		panic(fmt.Errorf("proto.Marshal MemberMeta err: %w", err))
	}
	return bs
}
func (it *Members) LocalState(join bool) []byte{
	return it.MemberStateManager.GetState(join)
}
func (it *Members) MergeRemoteState(buf []byte, join bool) {
	it.MemberStateManager.MergeState(buf, join)
}
func (it *Members) NotifyJoin(node *memberlist.Node) {

}
func (it *Members) NotifyLeave(node *memberlist.Node) {
	it.MemberStateManager.DeleteState(node.Name)
}
func (it *Members) NotifyUpdate(node *memberlist.Node) {

}
