package utils

import (
	"fmt"
	"sync"
	"time"
	"log"

	"github.com/lni/goutils/syncutil"
)

type ISyncClientSource[T any, K comparable] interface {
	List() ([]T, uint64, error)
	Subscribe(stopper *syncutil.Stopper, version uint64) (chan Stream[SyncStateView[T]], error)
	Key(item T) K
}

type ISyncClient[T any, K comparable] interface {
	Get(key K) (T, error)
	Close() error
}

type SyncClient[T any, K comparable] struct {
	src ISyncClientSource[T, K]

	mu sync.RWMutex
	dict map[K]T
	stopper *syncutil.Stopper
}

func NewSyncClient[T any, K comparable](src ISyncClientSource[T, K]) ISyncClient[T, K] {
	it := &SyncClient[T, K]{
		src: src,
		dict: make(map[K]T, 0),
		stopper: syncutil.NewStopper(),
	}
	it.stopper.RunWorker(it.updateWithRetry)
	return it
}

func (it *SyncClient[T, K]) Get(key K) (T, error) {
	it.mu.RLock()
	defer it.mu.RUnlock()
	item, ok := it.dict[key]
	if !ok {
		return *new(T), fmt.Errorf("Key NotFound: %v", key)
	}
	return item, nil
}


func (it *SyncClient[T, K]) updateWithRetry() {
	for {
		select {
		case <-it.stopper.ShouldStop():
			return
		default:
			if err := it.update(); err != nil {
				log.Println("[WARN]", fmt.Errorf("SyncClient.update err: %w", err))
				time.Sleep(time.Second)
				continue
			}
		}
	}
}

func (it *SyncClient[T, K]) update() error {
	items, version, err := it.src.List()
	if err != nil {
		return fmt.Errorf("src.List() %w", err)
	}
	it.mu.Lock()
	it.dict = make(map[K]T, len(items))
	for _, item := range items {
		key := it.src.Key(item)
		it.dict[key] = item
	}
	it.mu.Unlock()
	ch, err := it.src.Subscribe(it.stopper, version)
	if err != nil {
		return fmt.Errorf("src.Subscribe(%d) %w", version, err)
	}
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return fmt.Errorf("subscribe closed")
			}
			if ss.Error != nil {
				return fmt.Errorf("subscribe closed with err: %w", ss.Error)
			}
			it.mu.Lock()
			change := ss.Item
			key := it.src.Key(change.Item)
			switch change.Type {
			case SyncStateType_Add:
				it.dict[key] = change.Item
			case SyncStateType_Update:
				it.dict[key] = change.Item
			case SyncStateType_Remove:
				delete(it.dict, key)
			}
			it.mu.Unlock()
		case <-it.stopper.ShouldStop():
			return nil
		}
	}
}

func (it *SyncClient[T, K]) Close() error {
	it.stopper.Stop()
	return nil
}
