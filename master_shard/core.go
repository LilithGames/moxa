package master_shard

import (
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

type Core struct {
	state StateRoot

	shardIndex IIndexClient[string, *ShardSpec]
}

func NewCore() *Core {
	core := &Core{
		state: StateRoot{
			LastShardId: BaseSubShardID,
			ShardIndices: &Indices{},
		},
	}
	core.prepare()
	return core
}

func (it *Core) Healthz(req *HealthzRequest) (*HealthzResponse, error) {
	return &HealthzResponse{}, nil
}

func (it *Core) GetStateVersion(req *GetStateVersionRequest) (*GetStateVersionResponse, error) {
	return &GetStateVersionResponse{
		ShardsStateVersion: it.state.ShardsVersion,
		NodesStateVersion:  it.state.NodesVersion,
	}, nil
}

func (it *Core) prepare() {
	shardIndexSource := NewShardSpecIndexSource(it)
	shardIndex := NewDefaultIndex(it.state.ShardIndices)
	shardIndexer := NewFuncIndexer(1, func(spec *ShardSpec) []string {
		indices := []string{fmt.Sprintf("%s:%d", ShardSpecIndex_ShardID.String(), spec.ShardId)}
		for k, v := range spec.Labels {
			indices = append(indices, fmt.Sprintf("%s:%s", ShardSpecIndex_LabelKey.String(), k))
			indices = append(indices, fmt.Sprintf("%s:%s:%s", ShardSpecIndex_LabelKeyValue.String(), k, v))
		}
		return indices
	})
	it.shardIndex = NewIndexClient(shardIndexSource, shardIndex, shardIndexer)
}

func (it *Core) save(w io.Writer) error {
	data, err := proto.Marshal(&it.state)
	if err != nil {
		return fmt.Errorf("proto.Marshal err: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("io.Write err: %w", err)
	}
	return nil
}

func (it *Core) load(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll err: %w", err)
	}
	var state StateRoot
	if err := proto.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("proto.Unmarshal err: %w", err)
	}
	it.state = state
	it.prepare()
	return nil
}

func clone[T proto.Message](src T) T {
	return proto.Clone(src).(T)
}
