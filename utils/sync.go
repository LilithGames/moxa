package utils

import (
	"fmt"
)

type SyncState[T any] struct {
	Type SyncStateType
	Item T
}

type SyncStateView[T any] struct {
	Version uint64
	Type SyncStateType
	Item T
}

type ISync[T any] interface {
	ISyncReader[T]
	ISyncWriter[T]
}

type ISyncReader[T any] interface {
	Get(version uint64) ([]SyncStateView[T], error)
	Since(version uint64) ([]SyncStateView[T], error)
	Range(from uint64, to uint64) ([]SyncStateView[T], error)
}
type ISyncWriter[T any] interface {
	Push(version uint64, items ...SyncState[T]) error
}

type Sync[T any] struct {
	changes map[uint64][]SyncState[T]
	min     uint64
	curr    uint64
	size    int
}

func NewSync[T any](size int) ISync[T] {
	it := &Sync[T]{
		changes: make(map[uint64][]SyncState[T], size),
		min: 0,
		curr: 0,
		size: size,
	}
	it.shrink()
	return it
}

func (it *Sync[T]) shrink() {
	for {
		if it.min >= it.curr {
			return
		}
		if len(it.changes) <= it.size {
			return
		} 
		delete(it.changes, it.min)
		it.min++
	}
}

func (it *Sync[T]) Get(version uint64) ([]SyncStateView[T], error) {
	if version < it.min {
		return nil, fmt.Errorf("version %d expired, min: %d", version, it.min)
	}
	changes, ok := it.changes[version]
	if !ok {
		return nil, fmt.Errorf("version %d not found", version)
	}
	result := make([]SyncStateView[T], 0, len(changes))
	for _, change := range changes {
		result = append(result, SyncStateView[T]{Version: version, Type: change.Type, Item: change.Item})
	}
	return result, nil
}

func (it *Sync[T]) Since(version uint64) ([]SyncStateView[T], error) {
	if version < it.min {
		return nil, fmt.Errorf("version %d expired, min: %d", version, it.min)
	}
	result := make([]SyncStateView[T], 0)
	for i := version; i <= it.curr; i++ {
		if changes, ok := it.changes[i]; ok {
			for _, change := range changes {
				result = append(result, SyncStateView[T]{Version: i, Type: change.Type, Item: change.Item})
			}
		}
	}
	return result, nil
}

func (it *Sync[T]) Range(from uint64, to uint64) ([]SyncStateView[T], error) {
	if from > to {
		return nil, fmt.Errorf("from %d should less than to %d", from, to)
	}
	if from < it.min {
		return nil, fmt.Errorf("version %d expired, min: %d", from, it.min)
	}
	if to > it.curr {
		return nil, fmt.Errorf("version %d large than latest %d", to, it.curr)
	}
	result := make([]SyncStateView[T], 0)
	for i := from; i <= to; i++ {
		if changes, ok := it.changes[i]; ok {
			for _, change := range changes {
				result = append(result, SyncStateView[T]{Version: i, Type: change.Type, Item: change.Item})
			}
		}
	}
	return result, nil
}

func (it *Sync[T]) Push(version uint64, items ...SyncState[T]) error {
	if version <= it.curr {
		return fmt.Errorf("version(%d) should newer then curr(%d)", version, it.curr)
	}
	if _, ok := it.changes[version]; ok {
		return fmt.Errorf("version already exists (impossible cause)")
	}
	it.changes[version] = items
	it.curr = version
	it.shrink()
	return nil
}

