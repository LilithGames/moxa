package sub_shard

import (
	"encoding/binary"
	"io"
	"io/ioutil"

	sm "github.com/lni/dragonboat/v3/statemachine"
)

type ExampleStateMachine struct {
	ClusterID uint64
	NodeID    uint64
	Count     uint64
}

func NewExampleStateMachine(clusterID uint64, nodeID uint64) sm.IStateMachine {
	return &ExampleStateMachine{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
	}
}

func (it *ExampleStateMachine) Lookup(query interface{}) (interface{}, error) {
	return DragonboatSubShardLookup(it, query)
}

func (it *ExampleStateMachine) Update(data []byte) (sm.Result, error) {
	return DragonboatSubShardUpdate(it, data)
}

func (it *ExampleStateMachine) Get(req *GetRequest) (*GetResponse, error) {
	return &GetResponse{Value: it.Count}, nil
}

func (it *ExampleStateMachine) Incr(req *IncrRequest) (*IncrResponse, error) {
	it.Count++
	return &IncrResponse{Value: it.Count}, nil
}

func (it *ExampleStateMachine) SaveSnapshot(w io.Writer, fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, it.Count)
	_, err := w.Write(data)
	return err
}

func (it *ExampleStateMachine) RecoverFromSnapshot(r io.Reader, files []sm.SnapshotFile, done <-chan struct{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	it.Count = v
	return nil
}

func (it *ExampleStateMachine) Close() error { return nil }
