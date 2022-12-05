package master_shard

import (
	"fmt"
	"io"
	"sort"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/samber/lo"
	"github.com/serialx/hashring"
	"google.golang.org/protobuf/proto"
)

const BaseSubShardID uint64 = 10000000

type Core struct {
	state StateRoot
}

func NewCore() *Core {
	return &Core{
		state: StateRoot{
			LastShardId: BaseSubShardID,
		},
	}
}

func (it *Core) Healthz(req *HealthzRequest) (*HealthzResponse, error) {
	return &HealthzResponse{}, nil
}
func (it *Core) GetStateVersion(req *GetStateVersionRequest) (*GetStateVersionResponse, error) {
	return &GetStateVersionResponse{
		ShardsStateVersion: it.state.ShardsVersion,
		NodesStateVersion: it.state.NodesVersion,
	}, nil
}
func (it *Core) ListShards(req *ListShardsRequest) (*ListShardsResponse, error) {
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
func (it *Core) GetShard(req *GetShardRequest) (*GetShardResponse, error) {
	if shard, ok := it.state.Shards[req.Name]; ok {
		return &GetShardResponse{Shard: shard}, nil
	} else {
		return &GetShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
}
func (it *Core) CreateShard(req *CreateShardRequest) (*CreateShardResponse, error) {
	if _, ok := it.state.Shards[req.Name]; ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_AlreadyExists), fmt.Sprintf("Shard %s already exists", req.Name))
	}
	replica := int32(3)
	if req.Replica != nil {
		replica = *req.Replica
	}
	if replica == 0 {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "replica == 0")
	}
	shard := &ShardSpec{
		StateVersion: 0,
		ShardName:    req.Name,
		ShardId:      it.state.LastShardId,
		LastNodeId:   1,
	}
	nodes := it.prepareNodeCreateViews(req.Nodes)
	ring := hashring.New(lo.Keys(nodes))
	nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
	if !ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.Nodes))
	}
	create := lo.PickByKeys(nodes, nhids)
	it.updateNodes(shard, create, true)

	if it.state.Shards == nil {
		it.state.Shards = make(map[string]*ShardSpec, 0)
	}
	it.state.Shards[shard.ShardName] = shard
	it.state.LastShardId++
	it.state.ShardsVersion++
	return &CreateShardResponse{Shard: clone(shard)}, nil
}
func (it *Core) UpdateShard(req *UpdateShardRequest) (*UpdateShardResponse, error) {
	shard, ok := it.state.Shards[req.Name]
	if !ok {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
	replica := shard.Replica
	if req.Replica != nil {
		replica = *req.Replica
	}
	if replica == 0 {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "replica == 0")
	}

	nodes := it.prepareNodeCreateViews(req.Nodes)
	ring := hashring.New(lo.Keys(nodes))
	nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
	if !ok {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.Nodes))
	}
	apply := lo.PickByKeys(nodes, nhids)
	updated := it.updateNodes(shard, apply, false)

	if updated > 0 {
		it.state.ShardsVersion++
	}
	return &UpdateShardResponse{Updated: updated}, nil
}

func (it *Core) prepareNodeCreateViews(views []*NodeCreateView) map[string]*NodeCreateView {
	excludes := it.getExcludedNodeHostIds()
	views1 := lo.FilterMap(views, func(n *NodeCreateView, _ int) (*NodeCreateView, bool) {
		if _, ok := excludes[n.NodeHostId]; ok {
			return nil, false
		}
		return n, true
	})
	views2 := lo.SliceToMap(views1, func(n *NodeCreateView) (string, *NodeCreateView) {
		return n.NodeHostId, n
	})
	return views2
}

func (it *Core) getExcludedNodeHostIds() map[string]struct{} {
	nhids := make(map[string]struct{}, 0)
	for nhid, node := range it.state.Nodes {
		if v, ok := node.Labels["builtin.exclude"]; ok && v == "true" {
			nhids[nhid] = struct{}{}
		}
	}
	return nhids
}

func (it *Core) updateNodes(spec *ShardSpec, apply map[string]*NodeCreateView, initial bool) int32 {
	var updated int32 = 0
	curr := spec.Nodes
	nodeID := spec.LastNodeId
	result := make(map[string]*Node, len(apply))
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
	if updated > 0 {
		if initial {
			spec.Initials = result
		}
		spec.Nodes = result
		spec.Replica = int32(len(apply))
		spec.LastNodeId = nodeID
	}
	return updated
}

func (it *Core) DeleteShard(req *DeleteShardRequest) (*DeleteShardResponse, error) {
	if _, ok := it.state.Shards[req.Name]; ok {
		delete(it.state.Shards, req.Name)
		it.state.ShardsVersion++
		return &DeleteShardResponse{Updated: 1}, nil
	}
	return &DeleteShardResponse{Updated: 0}, nil
}

func (it *Core) GetNode(req *GetNodeRequest) (*GetNodeResponse, error) {
	node, ok := it.state.Nodes[req.NodeHostId]
	if !ok {
		return nil, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("node %s not found", req.NodeHostId))
	}
	return &GetNodeResponse{Node: clone(node)}, nil
}

func (it *Core) UpdateNode(req *UpdateNodeRequest) (*UpdateNodeResponse, error) {
	var updated uint64
	if it.state.Nodes == nil {
		it.state.Nodes = make(map[string]*NodeSpec, 0)
	}
	if prev, ok := it.state.Nodes[req.Node.NodeHostId]; ok {
		if prev.StateVersion != req.Node.StateVersion {
			return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("StateVersion %d should match with %d when updating", req.Node.StateVersion, prev.StateVersion))
		}
		if !proto.Equal(prev, req.Node) {
			node := clone(req.Node)
			node.StateVersion++
			it.state.Nodes[node.NodeHostId] = node
			updated++
		}
	} else {
		// create
		if req.Node.StateVersion != 0 {
			return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("StateVersion %d should be 0 when creating", req.Node.StateVersion))
		}
		node := clone(req.Node)
		node.StateVersion++
		it.state.Nodes[node.NodeHostId] = node
		updated++
	}
	if updated > 0 {
		it.state.NodesVersion++
	}
	return &UpdateNodeResponse{Updated: updated}, nil
}

func (it *Core) save(w io.Writer) error {
	data, err := proto.Marshal(&it.state)
	if err != nil {
		return fmt.Errorf("proto.Marshal err: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("io.Write err: %w", err)
	}
	return nil
}

func (it *Core) load(r io.Reader) error {
	data, err := io.ReadAll(r)
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

func clone[T proto.Message](src T) T {
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
