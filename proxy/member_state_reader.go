package proxy

import (
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/utils"
)

type MemberStateReader struct {
	syncClient utils.ISyncClient[*cluster.MemberState, string]
}

func NewMemberStateReader(syncClient utils.ISyncClient[*cluster.MemberState, string]) service.IMemberStateReader {
	return &MemberStateReader{syncClient: syncClient}
}

func (it *MemberStateReader) GetMemberStateList() (map[string]*cluster.MemberState, uint64) {
	return it.syncClient.List()
}

func (it *MemberStateReader) GetMemberStateVersion() uint64 {
	return it.syncClient.Version()
}

