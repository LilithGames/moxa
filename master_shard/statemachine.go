package master_shard

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/serialx/hashring"
	"google.golang.org/protobuf/proto"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/samber/lo"
)

const BaseSubShardID uint64 = 10000000

type StateMachine struct {
	ShardID uint64
	NodeID  uint64

	state StateRoot
}

func NewStateMachine(shardID uint64, nodeID uint64) sm.IStateMachine {
	return &StateMachine{
		ShardID: shardID,
		NodeID:  nodeID,
		state: StateRoot{
			LastShardId: BaseSubShardID,
		},
	}
}

func (it *StateMachine) Lookup(query interface{}) (interface{}, error) {
	return DragonboatMasterShardLookup(it, query)
}

func (it *StateMachine) Update(data []byte) (sm.Result, error) {
	return DragonboatMasterShardUpdate(it, data)
}

func (it *StateMachine) SaveSnapshot(w io.Writer, fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	data, err := proto.Marshal(&it.state)
	if err != nil {
		return fmt.Errorf("proto.Marshal err: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("io.Write err: %w", err)
	}
	return nil
}

func (it *StateMachine) RecoverFromSnapshot(r io.Reader, files []sm.SnapshotFile, done <-chan struct{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll err: %w", err)
	}
	var state StateRoot
	if err := proto.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("proto.Unmarshal err: %w", err)
	}
	it.state = state
	return nil
}

func (it *StateMachine) Close() error {
	return nil
}

func (it *StateMachine) Healthz(req *HealthzRequest) (*HealthzResponse, error) {
	return &HealthzResponse{}, nil
}
func (it *StateMachine) GetStateVersion(req *GetStateVersionRequest) (*GetStateVersionResponse, error) {
	return &GetStateVersionResponse{
		ShardsStateVersion: it.state.ShardsVersion,
	}, nil
}
func (it *StateMachine) ListShards(req *ListShardsRequest) (*ListShardsResponse, error) {
	shards := make([]*ShardSpec, 0, len(it.state.Shards))
	for _, shard := range it.state.Shards {
		// filter out shards not contains NodeHostId
		if req.NodeHostId != nil {
			if _, contains := shard.Nodes[*req.NodeHostId]; !contains {
				continue
			}
		}
		shards = append(shards, clone(shard))
	}
	return &ListShardsResponse{Shards: shards, StateVersion: it.state.ShardsVersion}, nil
}
func (it *StateMachine) GetShard(req *GetShardRequest) (*GetShardResponse, error) {
	if shard, ok := it.state.Shards[req.Name]; ok {
		return &GetShardResponse{Shard: shard}, nil
	} else {
		return &GetShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
}
func (it *StateMachine) CreateShard(req *CreateShardRequest) (*CreateShardResponse, error) {
	if _, ok := it.state.Shards[req.Name]; ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_AlreadyExists), fmt.Sprintf("Shard %s already exists", req.Name))
	}
	replica := int32(3)
	if req.Replica != nil {
		replica = *req.Replica
	}
	shard := ShardSpec{
		StateVersion:  0, 
		ShardName:  req.Name,
		ShardId:    it.state.LastShardId,
		Replica:    replica,
		Nodes:      make(map[string]*Node, 0),
	}
	nodes := lo.SliceToMap(req.Nodes, func(n *NodeCreateView) (string, *NodeCreateView) {
		return n.NodeHostId, n
	})
	ring := hashring.New(lo.Keys(nodes))
	nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
	if !ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.Nodes))
	}
	create := lo.PickByKeys(nodes, nhids)
	created, _ := it.updateNodes(shard.Nodes, create)

	shard.Nodes = created
	if it.state.Shards == nil {
		it.state.Shards = make(map[string]*ShardSpec, 0)
	}
	it.state.Shards[shard.ShardName] = &shard
	it.state.LastShardId++
	it.state.ShardsVersion++
	return &CreateShardResponse{Shard: clone(&shard)}, nil
}
func (it *StateMachine) UpdateShard(req *UpdateShardRequest) (*UpdateShardResponse, error) {
	shard, ok := it.state.Shards[req.Name]
	if !ok {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
	replica := shard.Replica
	if req.Replica != nil {
		replica = *req.Replica
	}
	nodes := lo.SliceToMap(req.Nodes, func(n *NodeCreateView) (string, *NodeCreateView) {
		return n.NodeHostId, n
	})
	ring := hashring.New(lo.Keys(nodes))
	nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
	if !ok {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.Nodes))
	}
	apply := lo.PickByKeys(nodes, nhids)
	applied, updated := it.updateNodes(shard.Nodes, apply)

	if updated > 0 {
		shard.Nodes = applied
		shard.Replica = replica
		it.state.ShardsVersion++
	}
	return &UpdateShardResponse{Updated: updated}, nil
}
func (it *StateMachine) updateNodes(curr map[string]*Node, apply map[string]*NodeCreateView) (map[string]*Node, int32) {
	var updated int32 = 0
	result := make(map[string]*Node, len(apply))
	nodeID := lastNodeID(curr)
	currKeys := lo.MapToSlice(curr, func(nhid string, _ *Node) string {
		return nhid
	})
	applyKeys := lo.MapToSlice(apply, func(nhid string, _ *NodeCreateView) string {
		return nhid
	})
	nhids := lo.Union(currKeys, applyKeys)
	sort.Strings(nhids)
	for _, nhid := range nhids {
		currNode, currOK := curr[nhid]
		applyNode, applyOK := apply[nhid]
		if currOK && applyOK {
			// addr cannot change for nodehost, so ignore applyNode.Addr
			result[nhid] = currNode
		} else if currOK && !applyOK {
			// delete
			updated++
			continue
		} else if !currOK && applyOK {
			result[nhid] = &Node{NodeHostId: applyNode.NodeHostId, NodeId: nodeID, Addr: applyNode.Addr}
			nodeID++
			updated++
		}
	}
	return result, updated
}

func (it *StateMachine) DeleteShard(req *DeleteShardRequest) (*DeleteShardResponse, error) {
	if _, ok := it.state.Shards[req.Name]; ok {
		delete(it.state.Shards, req.Name)
		it.state.ShardsVersion++
		return &DeleteShardResponse{Updated: 1}, nil
	}
	return &DeleteShardResponse{Updated: 0}, nil
}
func clone[T proto.Message](src T) T{
	return proto.Clone(src).(T)
}
func lastNodeID(nodes map[string]*Node) uint64 {
	var max uint64
	for _, node := range nodes {
		if node.NodeId > max {
			max = node.NodeId
		}
	}
	max++
	return max
}
