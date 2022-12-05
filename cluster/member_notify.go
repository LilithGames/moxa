package cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
	"google.golang.org/protobuf/proto"
	"github.com/hashicorp/go-multierror"
)

type NotifyOptions struct {
	ShardID *uint64
}

func (it *Members) Notify(notify *MemberNotify, opts NotifyOptions) error {
	bs, err := proto.Marshal(notify)
	if err != nil {
		return fmt.Errorf("proto.Marshal err: %w", err)
	}
	nodes := make([]*memberlist.Node, 0, it.Nums())
	it.Foreach(func(m MemberNode) bool {
		if opts.ShardID == nil {
			nodes = append(nodes, m.Node)
			return true
		}
		if m.State != nil {
			if _, ok := m.State.Shards[*opts.ShardID]; ok {
				nodes = append(nodes, m.Node)
				return true
			}
		}
		return true
	})
	var errs error
	for _, n := range nodes {
		if err := it.ml.SendBestEffort(n, bs); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (it *Members) NotifyMsg(bs []byte) {
	var notify MemberNotify
	if err := proto.Unmarshal(bs, &notify); err != nil {
		log.Println("[WARN]", fmt.Errorf("proto.Unmarhsal notify err: %w, throw away", err))
		return
	}
	switch notify.Name {
	case EventTopic_MigrationCancellingNotify.String():
		if notify.ShardId == nil {
			log.Println("[WARN]", fmt.Errorf("EventTopic_MigrationCancellingNotify missing shardID"))
			return
		}
		it.bus.PublishAsync(EventTopic_MigrationCancellingNotify.String(), &MigrationCancellingNotifyEvent{ShardId: *notify.ShardId})
	case EventTopic_MigrationRepairingNotify.String():
		if notify.ShardId == nil {
			log.Println("[WARN]", fmt.Errorf("EventTopic_MigrationRepairingNotify missing shardID"))
			return
		}
		it.bus.PublishAsync(EventTopic_MigrationRepairingNotify.String(), &MigrationRepairingNotifyEvent{ShardId: *notify.ShardId})
	default:
		log.Println("[WARN]", fmt.Errorf("unknown notify name %s", notify.Name))
	}
}
func (it *Members) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}
