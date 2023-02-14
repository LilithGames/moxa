package utils

import (
	"sync"
)

type Provider[T any] struct {
	item T
	mu sync.RWMutex
}

func NewProvider[T any](item T) *Provider[T] {
	return &Provider[T]{item: item}
}

func (it *Provider[T]) Get() T {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return it.item
}

func (it *Provider[T]) Set(item T) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.item = item
}
