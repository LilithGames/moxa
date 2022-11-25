package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/dragonboat/v3"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
)

type NodeHostService struct {
	cm cluster.Manager
	UnimplementedNodeHostServer
}

func RegisterNodeHostService(cs *ClusterService) {
	srv := &NodeHostService{cm: cs.ClusterManager()}
	RegisterNodeHostServer(cs.GrpcServiceRegistrar(), srv)
	cs.AddGrpcGatewayRegister(RegisterNodeHostHandlerFromEndpoint)
}

func (it *NodeHostService) Healthz(ctx context.Context, req *HealthzRequest) (*HealthzResponse, error) {
	if _, err := it.cm.Client().MasterShard().Healthz(ctx, &master_shard.HealthzRequest{}, runtime.WithClientTimeout(time.Second)); err != nil {
		return nil, fmt.Errorf("ShardClient.Healthz err: %w", err)
	}
	return &HealthzResponse{}, nil
}

func (it *NodeHostService) ShardHealthz(ctx context.Context, req *ShardHealthzRequest) (*ShardHealthzResponse, error) {
	if _, err := it.cm.Client().Raw(req.ShardId).Query(ctx, &runtime.DragonboatVoid{}, runtime.WithClientTimeout(time.Second)); err != nil {
		return &ShardHealthzResponse{}, err
	}
	return &ShardHealthzResponse{}, nil
}
func (it *NodeHostService) ShardNop(ctx context.Context, req *ShardNopRequest) (*ShardNopResponse, error) {
	if _, err := it.cm.Client().Raw(req.ShardId).Mutate(ctx, &runtime.DragonboatVoid{}); err != nil {
		return &ShardNopResponse{}, err
	}
	return &ShardNopResponse{}, nil
}

func (it *NodeHostService) Info(context.Context, *NodeHostInfoRequest) (*NodeHostInfoResponse, error) {
	nhi := it.cm.NodeHost().GetNodeHostInfo(dragonboat.NodeHostInfoOption{true})
	return &NodeHostInfoResponse{NodeHostInfo: NodeHostInfoSiri(nhi)}, nil
}

func (it *NodeHostService) AddNode(ctx context.Context, req *ShardAddNodeRequest) (*ShardAddNodeResponse, error) {
	m := it.getMember(req.NodeIndex)
	if m == nil {
		return nil, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
	}
	var nodeID uint64
	if req.NodeId != nil {
		nodeID = *req.NodeId
	} else if req.ShardId == cluster.MasterShardID {
		nodeID = m.MasterNodeId
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "nodeID required when adding subshard", req.NodeIndex)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	if err := it.cm.NodeHost().SyncRequestAddNode(ctx, req.ShardId, nodeID, m.RaftAddress, 0); err != nil {
		return nil, fmt.Errorf("SyncRequestAddNode err: %w", err)
	}
	return &ShardAddNodeResponse{}, nil
}
func (it *NodeHostService) RemoveNode(ctx context.Context, req *ShardRemoveNodeRequest) (*ShardRemoveNodeResponse, error) {
	m := it.getMember(req.NodeIndex)
	if m == nil {
		return &ShardRemoveNodeResponse{}, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
	}
	nodeID := it.getNodeID(req.ShardId, m.RaftAddress)
	if nodeID == nil {
		return &ShardRemoveNodeResponse{}, status.Errorf(codes.InvalidArgument, "NodeID for shard %d not found", req.ShardId)
	}
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()
	if err := it.cm.NodeHost().SyncRequestDeleteNode(ctx, req.ShardId, *nodeID, 0); err != nil {
		if errors.Is(err, dragonboat.ErrRejected) {
			return &ShardRemoveNodeResponse{}, nil
		}
		return nil, fmt.Errorf("SyncRequestDeleteNode err: %w", err)
	}
	return &ShardRemoveNodeResponse{Updated: 1}, nil
}

func (it *NodeHostService) ListNode(ctx context.Context, req *ShardListNodeRequest) (*ShardListNodeResponse, error) {
	nis := it.getNodeInfos(req.ShardId)
	return &ShardListNodeResponse{Nodes: nis}, nil
}

func (it *NodeHostService) TransferLeader(ctx context.Context, req *ShardTransferLeaderRequest) (*ShardTransferLeaderResponse, error) {
	m := it.getMember(req.NodeIndex)
	if m == nil {
		return nil, status.Errorf(codes.InvalidArgument, "NodeIndex %d not found", req.NodeIndex)
	}
	nodeID := it.getNodeID(req.ShardId, m.RaftAddress)
	if nodeID == nil {
		return &ShardTransferLeaderResponse{}, status.Errorf(codes.InvalidArgument, "NodeID for shard %d not found", req.ShardId)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	if err := it.cm.NodeHost().RequestLeaderTransfer(req.ShardId, *nodeID); err != nil {
		return nil, fmt.Errorf("RequestLeaderTransfer err: %w", err)
	}
	return &ShardTransferLeaderResponse{}, nil
}

func (it *NodeHostService) exportSnapshot(ctx context.Context, req *CreateSnapshotRequest) (*CreateSnapshotResponse, error) {
	if req.Leader && !it.cm.IsLeader(req.ShardId) {
		return nil, status.Errorf(codes.Unavailable, "Shard %d NodeHost NotLeader", req.ShardId)
	}
	expire := time.Minute*5
	version := it.cm.Version().DataVersion()
	nodes := it.cm.GetNodes(req.ShardId)
	if len(nodes) == 0 {
		return &CreateSnapshotResponse{}, status.Errorf(codes.InvalidArgument, "shardID %d not found", req.ShardId)
	}

	_, err := it.cm.Snapshot().GetLatest(req.ShardId, version, expire)
	if req.Force || err != nil {
		opt := dragonboat.SnapshotOption{
			Exported:   true,
			ExportPath: it.cm.Snapshot().Prepare(req.ShardId, version),
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
		s := it.cm.Snapshot().New(req.ShardId, version, index)
		if err := it.cm.Snapshot().Commit(s); err != nil {
			return nil, fmt.Errorf("Snapshot().Commit(%v), err: %w", s, err)
		}
	}
	snapshot, err := it.cm.Snapshot().GetLatest(req.ShardId, version, expire)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetLatest(%d) after SyncRequestSnapshot err: %v", req.ShardId, err)
	}
	if req.Leader && !it.cm.IsLeader(req.ShardId) {
		return nil, status.Errorf(codes.Unavailable, "Shard %d NodeHost NotLeader", req.ShardId)
	}
	return CreateSnapshotResponseSiri(snapshot, nodes), nil
}

func (it *NodeHostService) CreateSnapshot(ctx context.Context, req *CreateSnapshotRequest) (*CreateSnapshotResponse, error) {
	if req.Export {
		return it.exportSnapshot(ctx, req)
	}
	if req.Leader && !it.cm.IsLeader(req.ShardId) {
		return nil, status.Errorf(codes.Unavailable, "Shard %d NodeHost NotLeader", req.ShardId)
	}
	opt := dragonboat.SnapshotOption{}
	if req.CompactionOverhead != nil {
		opt.CompactionOverhead = *req.CompactionOverhead
		opt.OverrideCompactionOverhead = true
	}
	index, err := it.cm.NodeHost().SyncRequestSnapshot(ctx, req.ShardId, opt)
	if err != nil {
		if errors.Is(err, dragonboat.ErrRejected) {
			return nil, status.Errorf(codes.AlreadyExists, "SyncRequestSnapshot(%d) err: %v", err)
		} else if errors.Is(err, dragonboat.ErrSystemBusy) {
			return nil, status.Errorf(codes.Unavailable, "SyncRequestSnapshot(%d) err: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "SyncRequestSnapshot(%d) err: %v", err)
	}
	return &CreateSnapshotResponse{SnapshotIndex: index, Version: it.cm.Version().DataVersion()}, nil
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
			mshard := m.State.Shards[shardID]
			isLeader = mshard.IsLeader
		}
		return &NodeInfo{
			NodeId:     nodeID,
			NodeIndex:  m.Meta.NodeHostIndex,
			NodeHostId: m.Meta.NodeHostId,
			HostName:   m.Meta.HostName,
			Addr:       m.Meta.RaftAddress,
			Following:  true,
			Leading:    isLeader,
			Version:    m.Meta.CodeHash,
		}
	})
}
