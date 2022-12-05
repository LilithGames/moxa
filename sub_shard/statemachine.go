package sub_shard

import (
	"fmt"
	"io"

	sm "github.com/lni/dragonboat/v3/statemachine"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
)

type ExampleStateMachine struct {
	ClusterID uint64
	NodeID    uint64
	core      *Core
}

func NewExampleStateMachine(clusterID uint64, nodeID uint64) sm.IStateMachine {
	return &ExampleStateMachine{
		ClusterID: clusterID,
		NodeID:    nodeID,
		core:      NewCore(),
	}
}

func (it *ExampleStateMachine) Lookup(query interface{}) (interface{}, error) {
	return DragonboatSubShardLookup(it.core, query)
}

func (it *ExampleStateMachine) Update(data []byte) (sm.Result, error) {
	return DragonboatSubShardUpdate(it.core, data)
}

func (it *ExampleStateMachine) SaveSnapshot(w io.Writer, fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	if err := it.core.save(w); err != nil {
		return fmt.Errorf("core.save err: %w", err)
	}
	return nil
}

func (it *ExampleStateMachine) RecoverFromSnapshot(r io.Reader, files []sm.SnapshotFile, done <-chan struct{}) error {
	if err := it.core.load(r); err != nil {
		return fmt.Errorf("core.load err: %w", err)
	}
	return nil
}

func (it *ExampleStateMachine) SaveState(req *runtime.SaveStateRequest) (*runtime.SaveStateResponse, error) {
	if err := it.core.save(req.Writer); err != nil {
		return &runtime.SaveStateResponse{}, fmt.Errorf("core.save err: %w", err)
	}
	return &runtime.SaveStateResponse{}, nil
}

func (it *ExampleStateMachine) Migrate(req *runtime.MigrateRequest) (*runtime.MigrateResponse, error) {
	if err := it.core.load(req.Reader); err != nil {
		return &runtime.MigrateResponse{}, fmt.Errorf("core.load err: %w", err)
	}
	return &runtime.MigrateResponse{}, nil
}

func (it *ExampleStateMachine) Close() error {
	return nil
}
