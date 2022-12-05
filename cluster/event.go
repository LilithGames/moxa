package cluster

import (
	eventbus "github.com/LilithGames/go-event-bus/v4"
	"github.com/lni/dragonboat/v3/raftio"
)

type eventListener struct {
	bus *eventbus.EventBus
}

func (it *eventListener) LeaderUpdated(info raftio.LeaderInfo) {
	it.bus.PublishAsync(EventTopic_NodeHostLeaderUpdated.String(), &NodeHostLeaderUpdatedEvent{
		ShardId:  info.ClusterID,
		NodeId:   info.NodeID,
		Term:     info.Term,
		LeaderId: info.LeaderID,
	})
}
func (it *eventListener) NodeHostShuttingDown()             {}
func (it *eventListener) NodeUnloaded(info raftio.NodeInfo) {
	it.bus.PublishAsync(EventTopic_NodeHostNodeUnloaded.String(), &NodeHostNodeUnloadedEvent{
		ShardId: info.ClusterID,
		NodeId:  info.NodeID,
	})
}
func (it *eventListener) NodeReady(info raftio.NodeInfo) {
	it.bus.PublishAsync(EventTopic_NodeHostNodeReady.String(), &NodeHostNodeReadyEvent{
		ShardId: info.ClusterID,
		NodeId:  info.NodeID,
	})
}
func (it *eventListener) MembershipChanged(info raftio.NodeInfo) {
	it.bus.PublishAsync(EventTopic_NodeHostMembershipChanged.String(), &NodeHostMembershipChangedEvent{
		ShardId: info.ClusterID,
		NodeId:  info.NodeID,
	})
}
func (it *eventListener) ConnectionEstablished(info raftio.ConnectionInfo) {}
func (it *eventListener) ConnectionFailed(info raftio.ConnectionInfo)      {}
func (it *eventListener) SendSnapshotStarted(info raftio.SnapshotInfo)     {}
func (it *eventListener) SendSnapshotCompleted(info raftio.SnapshotInfo)   {}
func (it *eventListener) SendSnapshotAborted(info raftio.SnapshotInfo)     {}
func (it *eventListener) SnapshotReceived(info raftio.SnapshotInfo)        {}
func (it *eventListener) SnapshotRecovered(info raftio.SnapshotInfo)       {}
func (it *eventListener) SnapshotCreated(info raftio.SnapshotInfo)         {}
func (it *eventListener) SnapshotCompacted(info raftio.SnapshotInfo)       {}
func (it *eventListener) LogCompacted(info raftio.EntryInfo)               {}
func (it *eventListener) LogDBCompacted(info raftio.EntryInfo)             {}
