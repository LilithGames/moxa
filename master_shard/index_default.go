package master_shard

import (
	"fmt"
)

type DefaultIndex struct {
	index *Indices
}

func NewDefaultIndex(source *Indices) IIndex[string] {
	if source == nil {
		panic("source required")
	}
	return &DefaultIndex{index: source}
}

func (it *DefaultIndex) Set(key string, indices []string) {
	it.index.Meta[key] = &IndexMeta{Key: key, Indices: indices}
	for _, index := range indices {
		if data, ok := it.index.Indices[index]; ok {
			data.Keys = append(data.Keys, key)
		} else {
			it.index.Indices[index] = &IndexData{Index: index, Keys: []string{key}}
		}
	}
}
func (it *DefaultIndex) Del(key string) error {
	meta, ok := it.index.Meta[key]
	if !ok {
		return fmt.Errorf("Index meta not found: %s", key)
	}
	for _, index := range meta.Indices {
		delete(it.index.Indices, index)
	}
	delete(it.index.Meta, key)
	return nil
}
func (it *DefaultIndex) Get(index string) ([]string, error) {
	data, ok := it.index.Indices[index]
	if !ok {
		return nil, ErrIndexNotFound
	}
	return data.Keys, nil
}
func (it *DefaultIndex) IndexerVer() uint64 {
	return it.index.IndexerVersion
}
func (it *DefaultIndex) CommitIndexer(ver uint64) {
	it.index.IndexerVersion = ver
}
func (it *DefaultIndex) Clear() {
	it.index.Indices = make(map[string]*IndexData, 0)
	it.index.Meta = make(map[string]*IndexMeta, 0)
	it.index.IndexerVersion = 0
}
