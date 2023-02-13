package proxy

import (
	"context"
	"fmt"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/utils"
	"github.com/lni/goutils/syncutil"
)

type MemberSyncClientSource struct {
	client service.IClient
}

func NewMemberSyncClientSource(client service.IClient) utils.ISyncClientSource[*cluster.MemberState, string] {
	return &MemberSyncClientSource{client: client}
}

func (it *MemberSyncClientSource) List() ([]*cluster.MemberState, uint64, error) {
	resp, err := it.client.NodeHost().ListMemberState(context.TODO(), &service.ListMemberStateRequest{})
	if err != nil {
		return nil, 0, fmt.Errorf("client.NodeHost().ListMemberState err: %w", err)
	}
	return resp.Members, resp.Version, nil
}
func (it *MemberSyncClientSource) Subscribe(stopper *syncutil.Stopper, version uint64) (chan utils.Stream[utils.SyncStateView[*cluster.MemberState]], error) {
	ctx := utils.BindContext(stopper, context.TODO())
	session, err := it.client.NodeHost().SubscribeMemberState(ctx, &service.SubscribeMemberStateRequest{Version: version})
	if err != nil {
		return nil, fmt.Errorf("SubscribeMemberState err: %w", err)
	}
	ch := make(chan utils.Stream[utils.SyncStateView[*cluster.MemberState]], 0)
	stopper.RunWorker(func() {
		for {
			resp, err := session.Recv()
			if err != nil {
				ch <- utils.Stream[utils.SyncStateView[*cluster.MemberState]]{Error: err}
				close(ch)
				return
			}
			ssv := utils.SyncStateView[*cluster.MemberState]{Item: resp.State, Type: resp.Type, Version: resp.Version}
			ch <- utils.Stream[utils.SyncStateView[*cluster.MemberState]]{Item: ssv}
		}
	})
	return ch, nil
}
func (it *MemberSyncClientSource) Key(item *cluster.MemberState) string {
	return item.NodeHostId
}
