package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"reflect"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/utils"
)

type NodeHostService struct {
	cm cluster.Manager
	client IClient
	UnimplementedNodeHostServer
}

func RegisterNodeHostService(cs *ClusterService) {
	srv := &NodeHostService{cm: cs.ClusterManager()}
	RegisterNodeHostServer(cs.GrpcServiceRegistrar(), srv)
	cs.AddGrpcGatewayRegister(func(ctx context.Context, mux *gruntime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		return RegisterNodeHostHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", cs.Config().GrpcPort), opts)
	})
}

func (it *NodeHostService) Healthz(ctx context.Context, req *HealthzRequest) (*HealthzResponse, error) {
	if !it.cm.StartupReady().IsSet() {
		return nil, fmt.Errorf("startup not ready")
	}
	if _, err := it.cm.Client().MasterShard().Healthz(ctx, &master_shard.HealthzRequest{}, runtime.WithClientTimeout(time.Second)); err != nil {
		return nil, fmt.Errorf("ShardClient.Healthz err: %w", err)
	}
	return &HealthzResponse{}, nil
}

func (it *NodeHostService) ShardHealthz(ctx context.Context, req *ShardHealthzRequest) (*ShardHealthzResponse, error) {
	if err := it.cm.Client().Void(req.ShardId).VoidQuery(ctx, runtime.WithClientTimeout(time.Second)); err != nil {
		return &ShardHealthzResponse{}, fmt.Errorf("VoidQuery err: %w", err)
	}
	return &ShardHealthzResponse{}, nil
}
func (it *NodeHostService) ShardNop(ctx context.Context, req *ShardNopRequest) (*ShardNopResponse, error) {
	if err := it.cm.Client().Void(req.ShardId).VoidMutate(ctx); err != nil {
		return &ShardNopResponse{}, fmt.Errorf("VoidMutate err: %w", err)
	}
	return &ShardNopResponse{}, nil
}

func (it *NodeHostService) Info(context.Context, *NodeHostInfoRequest) (*NodeHostInfoResponse, error) {
	nhi := it.cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{false})
	return &NodeHostInfoResponse{NodeHostInfo: NodeHostInfoSiri(nhi)}, nil
}

func (it *NodeHostService) AddNode(ctx context.Context, req *ShardAddNodeRequest) (*ShardAddNodeResponse, error) {
	var addr string
	if req.Addr != nil {
		addr = *req.Addr
	} else if req.NodeIndex != nil {
		m := it.getMember(*req.NodeIndex)
		if m == nil {
			return nil, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
		}
		addr = m.RaftAddress
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "oneof Addr and NodeIndex required")
	}
	if err := it.cm.Client().Raw(req.ShardId).AddNode(ctx, req.NodeId, addr); err != nil {
		if errors.Is(err, dragonboat.ErrRejected) {
			return &ShardAddNodeResponse{}, nil
		}
		return nil, fmt.Errorf("Client.Raw.AddNode err: %w", err)
	}
	return &ShardAddNodeResponse{}, nil
}
func (it *NodeHostService) RemoveNode(ctx context.Context, req *ShardRemoveNodeRequest) (*ShardRemoveNodeResponse, error) {
	var nodeID uint64
	if req.NodeId != nil {
		nodeID = *req.NodeId
	} else if req.NodeIndex != nil {
		m := it.getMember(*req.NodeIndex)
		if m == nil {
			return &ShardRemoveNodeResponse{}, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
		}
		nid := it.getNodeID(req.ShardId, m.RaftAddress)
		if nid == nil {
			return &ShardRemoveNodeResponse{}, status.Errorf(codes.InvalidArgument, "NodeID for shard %d not found", req.ShardId)
		}
		nodeID = *nid
	}
	if err := it.cm.Client().Raw(req.ShardId).RemoveNode(ctx, nodeID); err != nil {
		return nil, fmt.Errorf("Client.Raw.RemoveNode err: %w", err)
	}
	return &ShardRemoveNodeResponse{Updated: 1}, nil
}

func (it *NodeHostService) SyncAddNode(ctx context.Context, req *SyncAddNodeRequest) (*SyncAddNodeResponse, error) {
	var addr string
	if req.Addr != nil {
		addr = *req.Addr
	} else if req.NodeIndex != nil {
		if meta := it.getMember(*req.NodeIndex); meta != nil {
			addr = meta.RaftAddress
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", *req.NodeIndex)
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Addr or NodeIndex required")
	}

	var updated uint64
	err := utils.RetryWithDelay(ctx, 5, time.Millisecond*10, func() (bool, error) {
		membership, err := it.cm.NodeHost().SyncGetClusterMembership(ctx, req.ShardId)
		if err != nil {
			return false, fmt.Errorf("NodeHost.SyncGetClusterMembership(%d) err: %w", req.ShardId, err)
		}
		if _, ok := membership.Nodes[req.NodeId]; ok {
			return false, nil
		}
		if req.Force == nil || !(*req.Force) {
			if !it.ShardFullHealthz(req.ShardId, membership) {
				return false, fmt.Errorf("shard %d not fully healthz", req.ShardId)
			}
		}
		if err := it.cm.NodeHost().SyncRequestAddNode(ctx, req.ShardId, req.NodeId, addr, membership.ConfigChangeID); err != nil {
			if errors.Is(err, dragonboat.ErrRejected) {
				return true, fmt.Errorf("SyncRequestAddNode(%v, %d) err: %w", req, membership.ConfigChangeID, err)
			}
			return false, fmt.Errorf("SyncRequestAddNode(%v, %d) err: %w", req, membership.ConfigChangeID, err)
		}
		updated = 1
		return false, nil
	})
	if err != nil {
		return &SyncAddNodeResponse{}, fmt.Errorf("RetryWithDelay err: %w", err)
	}
	return &SyncAddNodeResponse{Updated: updated}, nil
}

func (it *NodeHostService) SyncRemoveNode(ctx context.Context, req *SyncRemoveNodeRequest) (*SyncRemoveNodeResponse, error) {
	var nodeID uint64
	if req.NodeId != nil {
		nodeID = *req.NodeId
	} else if req.NodeIndex != nil {
		if meta := it.getMember(*req.NodeIndex); meta != nil {
			if nid := it.getNodeID(req.ShardId, meta.RaftAddress); nid != nil {
				nodeID = *nid
			} else {
				return nil, status.Errorf(codes.InvalidArgument, "nodeID of NodeIndex %d not found", *req.NodeIndex)
			}
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", *req.NodeIndex)
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "NodeId or NodeIndex required")
	}

	var updated uint64
	err := utils.RetryWithDelay(ctx, 5, time.Millisecond*10, func() (bool, error) {
		membership, err := it.cm.NodeHost().SyncGetClusterMembership(ctx, req.ShardId)
		if err != nil {
			return false, fmt.Errorf("NodeHost.SyncGetClusterMembership(%d) err: %w", req.ShardId, err)
		}
		if _, ok := membership.Nodes[nodeID]; !ok {
			return false, nil
		}
		if req.Force == nil || !(*req.Force) {
			if !it.ShardFullHealthz(req.ShardId, membership) {
				return false, fmt.Errorf("shard %d not fully healthz", req.ShardId)
			}
		}
		if err := it.cm.NodeHost().SyncRequestDeleteNode(ctx, req.ShardId, nodeID, membership.ConfigChangeID); err != nil {
			if errors.Is(err, dragonboat.ErrRejected) {
				return true, fmt.Errorf("SyncRequestDeleteNode(%v, %d) err: %w", req, membership.ConfigChangeID, err)
			}
			return false, fmt.Errorf("SyncRequestDeleteNode(%v, %d) err: %w", req, membership.ConfigChangeID, err)
		}
		updated = 1
		return false, nil
	})
	if err != nil {
		return &SyncRemoveNodeResponse{}, fmt.Errorf("RetryWithDelay err: %w", err)
	}
	return &SyncRemoveNodeResponse{Updated: updated}, nil

}

func (it *NodeHostService) ShardFullHealthz(shardID uint64, membership *dragonboat.Membership) bool {
	spec := membership.Nodes
	status := make(map[uint64]string, 0)
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		if m.State != nil {
			if shard, ok := m.State.Shards[shardID]; ok {
				if reflect.DeepEqual(spec, shard.Nodes) {
					status[shard.NodeId] = m.Meta.RaftAddress
				} else {
					return false
				}
			}
		}
		return true
	})
	return reflect.DeepEqual(spec, status)

}

func (it *NodeHostService) ListNode(ctx context.Context, req *ShardListNodeRequest) (*ShardListNodeResponse, error) {
	nis := it.getNodeInfos(req.ShardId)
	return &ShardListNodeResponse{Nodes: nis}, nil
}

func (it *NodeHostService) TransferLeader(ctx context.Context, req *ShardTransferLeaderRequest) (*ShardTransferLeaderResponse, error) {
	var nodeID uint64
	if req.NodeId != nil {
		nodeID = *req.NodeId
	} else if req.NodeIndex != nil {
		m := it.getMember(*req.NodeIndex)
		if m == nil {
			return &ShardTransferLeaderResponse{}, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
		}
		nid := it.getNodeID(req.ShardId, m.RaftAddress)
		if nid == nil {
			return &ShardTransferLeaderResponse{}, status.Errorf(codes.InvalidArgument, "NodeID for shard %d not found", req.ShardId)
		}
		nodeID = *nid
	}
	if err := it.cm.NodeHost().RequestLeaderTransfer(req.ShardId, nodeID); err != nil {
		return nil, fmt.Errorf("RequestLeaderTransfer err: %w", err)
	}
	return &ShardTransferLeaderResponse{}, nil
}

func (it *NodeHostService) CreateSnapshot(ctx context.Context, req *CreateSnapshotRequest) (*CreateSnapshotResponse, error) {
	opt := dragonboat.SnapshotOption{}
	if req.CompactionOverhead != nil {
		opt.CompactionOverhead = *req.CompactionOverhead
		opt.OverrideCompactionOverhead = true
	}
	if req.ExportPath != nil {
		opt.Exported = true
		opt.ExportPath = *req.ExportPath
	}
	index, err := it.cm.NodeHost().SyncRequestSnapshot(ctx, req.ShardId, opt)
	if err != nil {
		if errors.Is(err, dragonboat.ErrRejected) {
			return nil, status.Errorf(codes.AlreadyExists, "SyncRequestSnapshot(%d) err: %v", req.ShardId, err)
		} else if errors.Is(err, dragonboat.ErrSystemBusy) {
			return nil, status.Errorf(codes.Unavailable, "SyncRequestSnapshot(%d) err: %v", req.ShardId, err)
		}
		return nil, status.Errorf(codes.Internal, "SyncRequestSnapshot(%d) err: %v", req.ShardId, err)
	}
	return &CreateSnapshotResponse{SnapshotIndex: index}, nil
}

func (it *NodeHostService) ListMemberState(ctx context.Context, req *ListMemberStateRequest) (*ListMemberStateResponse, error) {
	state, version := it.cm.Members().GetMemberStateList()
	return &ListMemberStateResponse{Members: lo.Values(state), Version: version}, nil
}
func (it *NodeHostService) SubscribeMemberState(req *SubscribeMemberStateRequest, svr NodeHost_SubscribeMemberStateServer) error {
	stream, done := it.cm.Members().ChangesSince(req.Version)
	for {
		select {
		case ss, ok := <-stream:
			if !ok {
				return fmt.Errorf("stream eof")
			}
			if ss.Error != nil {
				return fmt.Errorf("stream err: %w", ss.Error)
			}
			change := ss.Item
			if err := svr.Send(&SubscribeMemberStateResponse{State: change.Item, Type: change.Type, Version: change.Version}); err != nil {
				return fmt.Errorf("svr.Send() err: %w", err)
			}
		case <-svr.Context().Done():
			close(done)
			return nil
		}
	}
}

func (it *NodeHostService) getMember(nhIndex int32) *cluster.MemberMeta {
	var member *cluster.MemberMeta
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		if m.Meta.NodeHostIndex == nhIndex {
			member = m.Meta
			return false
		}
		return true
	})
	return member
}

func (it *NodeHostService) getNodeID(shardID uint64, addr string) *uint64 {
	nhi := it.cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	for _, shard := range nhi.ClusterInfoList {
		if shard.ClusterID == shardID {
			for nodeID, nodeAddr := range shard.Nodes {
				if nodeAddr == addr {
					return &nodeID
				}
			}
		}
	}
	return nil
}

func (it *NodeHostService) getNodeInfos(shardID uint64) []*NodeInfo {
	members := make(map[string]cluster.MemberNode)
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		members[m.Meta.RaftAddress] = m
		return true
	})
	nodes := it.cm.GetNodes(shardID)
	return lo.MapToSlice(nodes, func(nodeID uint64, addr string) *NodeInfo {
		m := members[addr]
		isLeader := false
		if m.State != nil {
			if mshard, ok := m.State.Shards[shardID]; ok {
				isLeader = mshard.IsLeader
			}
		}
		return &NodeInfo{
			NodeId:     nodeID,
			NodeIndex:  m.Meta.NodeHostIndex,
			NodeHostId: m.Meta.NodeHostId,
			HostName:   m.Meta.HostName,
			Addr:       m.Meta.RaftAddress,
			Following:  true,
			Leading:    isLeader,
		}
	})
}
