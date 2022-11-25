package cluster

import (
	"sync"
	"fmt"

	"google.golang.org/protobuf/proto"
	"github.com/samber/lo"
)

type IMemberStateManager interface {
	SetMemberState(s MemberState) error
	GetMemberState(nhid string) *MemberState
	GetMemberStateList() map[string]*MemberState
}
type MemberStateManager struct {
	mu    sync.RWMutex

	nhid  string
	state *MemberGlobalState
}

func NewMemberStateManager(nodeHostID string) *MemberStateManager {
	return &MemberStateManager{
		mu:   sync.RWMutex{},
		nhid: nodeHostID,
		state: &MemberGlobalState{
			Members: make(map[string]*MemberState, 0),
		},
	}
}

func (it *MemberStateManager) SetMemberState(s MemberState) error {
	if s.NodeHostId != it.nhid {
		return fmt.Errorf("cannot set other nodehost")
	}

	it.mu.Lock()
	defer it.mu.Unlock()
	if it.state.Members == nil {
		it.state.Members = make(map[string]*MemberState, 0)
	}
	it.state.Members[it.nhid] = &s
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
func (it *MemberStateManager) GetMemberStateList() map[string]*MemberState {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return proto.Clone(it.state).(*MemberGlobalState).Members
}

func (it *MemberStateManager) GetState(join bool) []byte {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return lo.Must1(proto.Marshal(it.state))
	
}
func (it *MemberStateManager) MergeState(buf []byte, join bool) {
	rstate := &MemberGlobalState{}
	lo.Must0(proto.Unmarshal(buf, rstate))
	if self, ok := it.state.Members[it.nhid]; ok {
		if rstate.Members == nil {
			rstate.Members = make(map[string]*MemberState, 0)
		}
		rstate.Members[it.nhid] = self
	} else {
		delete(rstate.Members, it.nhid)
	}
	it.mu.Lock()
	defer it.mu.Unlock()
	it.state = rstate
}
