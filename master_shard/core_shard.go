package master_shard

import (
	"fmt"
	"sort"
	"reflect"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/antonmedv/expr"
	"github.com/samber/lo"
	"github.com/serialx/hashring"
)

const BaseSubShardID uint64 = 10000000

func (it *Core) ListShards(req *ListShardsRequest) (*ListShardsResponse, error) {
	if req.Query != nil && req.IndexQuery != nil {
		return &ListShardsResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("both Query and IndexQuery provided"))
	}
	if req.Query != nil {
		shards, err := it.listShardsByQuery(*req.Query)
		if err != nil {
			return &ListShardsResponse{}, err
		}
		return &ListShardsResponse{Shards: shards, StateVersion: it.state.ShardsVersion}, nil
	} else if req.IndexQuery != nil {
		shards, err := it.listShardByIndex(req.IndexQuery)
		if err != nil {
			return &ListShardsResponse{}, err
		}
		return &ListShardsResponse{Shards: shards, StateVersion: it.state.ShardsVersion}, nil
	} else {
		shards := it.listShards()
		return &ListShardsResponse{Shards: shards, StateVersion: it.state.ShardsVersion}, nil
	}
	return &ListShardsResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, "impossible case")
}
func (it *Core) listShards() []*ShardSpec {
	shards := make([]*ShardSpec, 0, len(it.state.Shards))
	for _, shard := range it.state.Shards {
		shards = append(shards, clone(shard))
	}
	return shards
}
func (it *Core) listShardsByQuery(query string) ([]*ShardSpec, error) {
	predict, err := expr.Compile(query)
	if err != nil {
		return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("invalid query expression: %v", err))
	}
	shards := make([]*ShardSpec, 0, len(it.state.Shards))
	for _, shard := range it.state.Shards {
		v, err := expr.Run(predict, map[string]any{"it": shard})
		if err != nil {
			return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("runtime invalid query expression: %v", err))
		}
		include, ok := v.(bool)
		if !ok {
			return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("invalid query expression return type: %T", v))
		}
		if !include {
			continue
		}
		shards = append(shards, clone(shard))
	}
	return shards, nil

}
func (it *Core) listShardByIndex(query *ShardIndexQuery) ([]*ShardSpec, error) {
	shards, err := it.shardIndex.Get(fmt.Sprintf("%s:%s", query.Name, query.Value))
	if err != nil {
		return nil, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, fmt.Sprintf("index.Get(%v) err: %v", query, err))
	}
	return shards, nil
}

func (it *Core) GetShard(req *GetShardRequest) (*GetShardResponse, error) {
	if shard, ok := it.state.Shards[req.Name]; ok {
		return &GetShardResponse{Shard: shard}, nil
	} else {
		return &GetShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
}
func (it *Core) CreateShard(req *CreateShardRequest) (*CreateShardResponse, error) {
	if req.Name == "" {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "Name required")
	}
	if _, ok := it.state.Shards[req.Name]; ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_AlreadyExists), fmt.Sprintf("Shard %s already exists", req.Name))
	}
	replica := int32(3)
	if req.NodesView.Replica != nil {
		replica = *req.NodesView.Replica
	}
	if replica == 0 {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "replica == 0")
	}
	if req.ProfileName == "" {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "ProfileName required")
	}
	shard := &ShardSpec{
		StateVersion: 1,
		ShardName:    req.Name,
		ShardId:      it.state.LastShardId,
		ProfileName:  req.ProfileName,
		LastNodeId:   1,
		Labels:       req.Labels,
	}
	nodes := it.prepareNodeCreateViews(req.NodesView.Nodes)
	ring := hashring.New(lo.Keys(nodes))
	nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
	if !ok {
		return &CreateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.NodesView.Nodes))
	}
	create := lo.PickByKeys(nodes, nhids)
	it.updateNodes(shard, create, true)

	if it.state.Shards == nil {
		it.state.Shards = make(map[string]*ShardSpec, 0)
	}
	it.state.Shards[shard.ShardName] = shard
	it.state.LastShardId++
	it.state.ShardsVersion++
	it.shardIndex.Update(shard)
	return &CreateShardResponse{Shard: clone(shard)}, nil
}
func (it *Core) UpdateShard(req *UpdateShardRequest) (*UpdateShardResponse, error) {
	shard, ok := it.state.Shards[req.Name]
	if !ok {
		return &UpdateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_NotFound), fmt.Sprintf("Shard %s not found", req.Name))
	}
	var updated int32
	shard = clone(shard)

	if req.NodesView != nil {
		replica := shard.Replica
		if req.NodesView.Replica != nil {
			replica = *req.NodesView.Replica
		}
		if replica == 0 {
			return &UpdateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, "replica == 0")
		}

		nodes := it.prepareNodeCreateViews(req.NodesView.Nodes)
		ring := hashring.New(lo.Keys(nodes))
		nhids, ok := ring.GetNodes(fmt.Sprintf("shard-%d", shard.ShardId), int(replica))
		if !ok {
			return &UpdateShardResponse{}, runtime.NewDragonboatError(runtime.ErrCodeInternal, fmt.Sprintf("hashring GetNodes(%v) fail", req.NodesView.Nodes))
		}
		apply := lo.PickByKeys(nodes, nhids)
		updated += it.updateNodes(shard, apply, false)
	}
	if req.UpdateView != nil {
		if req.UpdateView.StateVersion != nil {
			if *req.UpdateView.StateVersion != shard.StateVersion {
				return &UpdateShardResponse{}, runtime.NewDragonboatError(int32(ErrCode_VersionExpired), fmt.Sprintf("version expired: %d != %d", *req.UpdateView.StateVersion, shard.StateVersion))
			}
		}
		if !reflect.DeepEqual(shard.Labels, req.UpdateView.Labels) {
			shard.Labels = req.UpdateView.Labels
			updated++
		}
	}

	if updated > 0 {
		shard.StateVersion++
		it.state.Shards[shard.ShardName] = shard
		it.state.ShardsVersion++
		it.shardIndex.Update(shard)
	}
	return &UpdateShardResponse{Updated: updated}, nil
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
		it.shardIndex.Delete(req.Name)
		return &DeleteShardResponse{Updated: 1}, nil
	}
	return &DeleteShardResponse{Updated: 0}, nil
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
