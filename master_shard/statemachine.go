package master_shard

import (
	"io"
	"fmt"

	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
)

type StateMachine struct {
	ShardID uint64
	NodeID  uint64

	core *Core
}

func NewStateMachine(shardID uint64, nodeID uint64) sm.IStateMachine {
	return &StateMachine{
		ShardID: shardID,
		NodeID:  nodeID,
		core:    NewCore(),
	}
}

func (it *StateMachine) Lookup(query interface{}) (interface{}, error) {
	return DragonboatMasterShardLookup(it.core, query)
}

func (it *StateMachine) Update(data []byte) (sm.Result, error) {
	return DragonboatMasterShardUpdate(it.core, data)
}

func (it *StateMachine) SaveSnapshot(w io.Writer, fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	return it.core.save(w)
}

func (it *StateMachine) RecoverFromSnapshot(r io.Reader, files []sm.SnapshotFile, done <-chan struct{}) error {
	return it.core.load(r)
}

func (it *StateMachine) SaveState(req *runtime.SaveStateRequest) (*runtime.SaveStateResponse, error) {
	if err := it.core.save(req.Writer); err != nil {
		return &runtime.SaveStateResponse{}, fmt.Errorf("core.save err: %w", err)
	}
	return &runtime.SaveStateResponse{}, nil
}

func (it *StateMachine) Migrate(req *runtime.MigrateRequest) (*runtime.MigrateResponse, error) {
	if err := it.core.load(req.Reader); err != nil {
		return &runtime.MigrateResponse{}, fmt.Errorf("core.load err: %w", err)
	}
	return &runtime.MigrateResponse{}, nil
}


func (it *StateMachine) Close() error {
	return nil
}
