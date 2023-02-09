package utils

import (
	"fmt"
	"sync"
)

type Stream[T any] struct {
	Item T
	Error error
}

type StreamWriter[T any] struct {
	s chan Stream[T]
	once sync.Once
}

func NewStreamWriter[T any](s chan Stream[T]) *StreamWriter[T] {
	return &StreamWriter[T]{s: s}
}

func (it *StreamWriter[T]) Done() {
	it.once.Do(func() {
		close(it.s)
	})
}

func (it *StreamWriter[T]) PutItem(item T) {
	it.s <-Stream[T]{Item: item}
}
func (it *StreamWriter[T]) NonblockingPutItem(item T) error {
	select {
	case it.s <-Stream[T]{Item: item}:
		return nil
	default:
		return fmt.Errorf("stream full")
	}
}
func (it *StreamWriter[T]) PutError(err error) {
	select {
	case it.s <-Stream[T]{Error: err}:
		it.Done()
	default:
		go func() {
			it.s <-Stream[T]{Error: err}
			it.Done()
		}()
	}
}

