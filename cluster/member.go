package cluster

import (
	"fmt"
	"time"
	"log"

	"github.com/hashicorp/memberlist"
	"google.golang.org/protobuf/proto"
	"github.com/lni/dragonboat/v3"
	// "github.com/lni/goutils/syncutil"
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
	Nums() int
	Stop() error

}

type Members struct {
	seed string

	meta *MemberMeta
	ml   *memberlist.Memberlist
	*MemberStateManager
}

func NewMembers(meta MemberMeta) (IMembers, error) {
	it := &Members{
		seed: "moxa-headless:7946",
		meta: &meta,
		MemberStateManager: NewMemberStateManager(meta.NodeHostId),
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
	if _, err := it.ml.Join([]string{it.seed}); err != nil {
		return fmt.Errorf("memberlist.Join err %w", err)
	}
	return nil
}

func (it *Members) UpdateNodeHostInfo(nhi *dragonboat.NodeHostInfo) error {
	shards := make(map[uint64]*MemberShard)
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
	if err := it.SetMemberState(MemberState{NodeHostId: nhi.NodeHostID, Shards: shards}); err != nil {
		return fmt.Errorf("SetMemberState err: %w", err)
	}
	return nil
}

func (it *Members) Foreach(fn func(MemberNode) bool) {
	mstate := it.GetMemberStateList()
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

func (it *Members) Nums() int {
	return it.ml.NumMembers()
}

func (it *Members) Stop() error {
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

func (it *Members) NotifyMsg([]byte) {}
func (it *Members) GetBroadcasts(overhead, limit int) [][]byte { return nil }
