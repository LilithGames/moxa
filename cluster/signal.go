package cluster

import (
	"sync"
)

type Signal struct {
	resettable bool
	ch *chan struct{}
	once *sync.Once
	mu   sync.RWMutex
}

func NewSignal(resettable bool) *Signal {
	ch := make(chan struct{}, 0)
	return &Signal{
		resettable: resettable,
		ch: &ch,
		once: &sync.Once{},
	}
}

func (it *Signal) Set() {
	if it.resettable {
		it.mu.RLock()
		defer it.mu.RUnlock()
	}
	it.once.Do(func(){
		close(*it.ch)
	})
}

func (it *Signal) Reset() {
	if !it.resettable {
		panic("cannot reset signal")
	}
	it.mu.Lock()
	defer it.mu.Unlock()
	it.once = &sync.Once{}
	ch := make(chan struct{}, 0)
	it.ch = &ch
}

func (it *Signal) IsSet() bool {
	if it.resettable {
		it.mu.RLock()
		defer it.mu.RUnlock()
	}
	select {
	case _, ok := <-*it.ch:
		if !ok {
			return true
		}
		return false
	default:
		return false
	}
}

func (it *Signal) Setted() <-chan struct{} {
	if it.resettable {
		it.mu.RLock()
		defer it.mu.RUnlock()
	}
	return *it.ch
}
