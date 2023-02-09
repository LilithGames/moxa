package cluster

import (
	"fmt"
	"sync"

	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	"github.com/LilithGames/moxa/utils"
)

type IMemberStateManager interface {
	SetMemberState(s MemberState) error
	GetMemberState(nhid string) *MemberState
	GetMemberStateList() (map[string]*MemberState, uint64)
	GetMemberStateVersion() uint64
	utils.ISyncReader[*MemberState]
}
type MemberStateManager struct {
	mu sync.RWMutex

	nhid  string
	bus   *eventbus.EventBus
	state *MemberGlobalState
	utils.ISync[*MemberState]
}

func NewMemberStateManager(nodeHostID string, bus *eventbus.EventBus) *MemberStateManager {
	return &MemberStateManager{
		mu:   sync.RWMutex{},
		bus:  bus,
		nhid: nodeHostID,
		state: &MemberGlobalState{
			Members: make(map[string]*MemberState, 0),
		},
		ISync: utils.NewSync[*MemberState](1024),
	}
}

func (it *MemberStateManager) SetMemberState(s MemberState) error {
	if s.NodeHostId != it.nhid {
		return fmt.Errorf("cannot set other nodehost")
	}
    it.updateMemberState(&s)
	return nil
}

func (it *MemberStateManager) GetMemberState(nhid string) *MemberState {
	it.mu.RLock()
	defer it.mu.RUnlock()
	if s, ok := it.state.Members[it.nhid]; ok {
		return proto.Clone(s).(*MemberState)
	} else {
		return nil
	}
}
func (it *MemberStateManager) GetMemberStateList() (map[string]*MemberState, uint64) {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return proto.Clone(it.state).(*MemberGlobalState).Members, it.state.Version
}

func (it *MemberStateManager) GetMemberStateVersion() uint64 {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return it.state.Version
}

func (it *MemberStateManager) GetState(join bool) []byte {
	it.mu.RLock()
	defer it.mu.RUnlock()
	self := it.state.Members[it.nhid]
	return lo.Must1(proto.Marshal(self))
}

func (it *MemberStateManager) MergeState(buf []byte, join bool) {
	mstate := &MemberState{}
	lo.Must0(proto.Unmarshal(buf, mstate))
	if mstate.NodeHostId != it.nhid {
		it.updateMemberState(mstate)
	}
}

func (it *MemberStateManager) DeleteState(nhid string) {
	it.mu.Lock()
	defer it.mu.Unlock()
	if mstate, ok := it.state.Members[nhid]; ok {
		delete(it.state.Members, nhid)
		it.ISync.Push(it.state.Version, utils.SyncState[*MemberState]{Type: utils.SyncStateType_Remove, Item: mstate})
		it.bus.PublishAsync(EventTopic_MemberStateChanged.String(), &MemberStateChangedEvent{Version: it.state.Version, NodeHostId: mstate.NodeHostId})
		it.state.Version++
	}
}

func (it *MemberStateManager) updateMemberState(mstate *MemberState) {
	it.mu.Lock()
	defer it.mu.Unlock()
	if it.state.Members == nil {
		it.state.Members = make(map[string]*MemberState, 0)
	}
	if pstate, ok := it.state.Members[mstate.NodeHostId]; ok {
		if !proto.Equal(pstate, mstate) {
			it.state.Members[mstate.NodeHostId] = mstate
			it.ISync.Push(it.state.Version, utils.SyncState[*MemberState]{Type: utils.SyncStateType_Update, Item: mstate})
			it.bus.PublishAsync(EventTopic_MemberStateChanged.String(), &MemberStateChangedEvent{Version: it.state.Version, NodeHostId: mstate.NodeHostId})
			it.state.Version++
		}
	} else {
		it.state.Members[mstate.NodeHostId] = mstate
		it.ISync.Push(it.state.Version, utils.SyncState[*MemberState]{Type: utils.SyncStateType_Add, Item: mstate})
		it.bus.PublishAsync(EventTopic_MemberStateChanged.String(), &MemberStateChangedEvent{Version: it.state.Version, NodeHostId: mstate.NodeHostId})
		it.state.Version++
	}
}
