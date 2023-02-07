package master_shard

import (
	"fmt"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

const LabelBuiltinExclude string = "builtin.exclude"

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
		if v, ok := node.Labels[LabelBuiltinExclude]; ok && v == "true" {
			nhids[nhid] = struct{}{}
		}
	}
	return nhids
}
