package master_shard

import (
	"fmt"
	"errors"
)

var ErrIndexNotFound = errors.New("ErrIndexNotFound")

type IIndex[K comparable] interface {
	Set(key K, indices []string)
	Del(key K) error
	Get(index string) ([]K, error)

	IndexerVer() uint64
	CommitIndexer(ver uint64)

	Clear()
}

type Indexer[T any] interface {
	Index(item T) []string
	Ver() uint64
}

type IIndexSource[K comparable, T any] interface {
	List() []T
	Get(key K) (T, error)
	Key(item T) K
}

type IIndexClient[K comparable, T any] interface {
	Get(index string) ([]T, error)
	Update(item T)
	Delete(key K)
}

type IndexClient[K comparable, T any] struct {
	src IIndexSource[K, T]
	index IIndex[K]
	indexer Indexer[T]
}

func NewIndexClient[K comparable, T any](src IIndexSource[K, T], index IIndex[K], indexer Indexer[T]) IIndexClient[K, T] {
	it := &IndexClient[K, T]{src: src, index: index, indexer: indexer}
	if it.indexer.Ver() != it.index.IndexerVer() {
		it.reset()
	}
	return it
}

func (it *IndexClient[K, T]) Get(index string) ([]T, error) {
	keys, err := it.index.Get(index)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, index)
	}
	items := make([]T, 0, len(keys))
	for _, key := range keys {
		item, err := it.src.Get(key)
		if err != nil {
			return nil, fmt.Errorf("src.Get(%v) err: %w", key, err)
		}
		items = append(items, item)
	}
	return items, nil
}
func (it *IndexClient[K, T]) Update(item T) {
	it.updateItem(item)
	return
}
func (it *IndexClient[K, T]) Delete(key K) {
	it.deleteItem(key)
	return 
}

func (it *IndexClient[K, T]) updateItem(item T) {
	key := it.src.Key(item)
	_ = it.index.Del(key)
	indices := it.indexer.Index(item)
	it.index.Set(key, indices)
}

func (it *IndexClient[K, T]) deleteItem(key K) {
	it.index.Del(key)
}

func (it *IndexClient[K, T]) reset() {
	it.index.Clear()
	items := it.src.List()
	for _, item := range items {
		key := it.src.Key(item)
		indices := it.indexer.Index(item)
		it.index.Set(key, indices)
	}
	it.index.CommitIndexer(it.indexer.Ver())
}
