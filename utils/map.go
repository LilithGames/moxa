package utils

import (
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
)

func ParallelMap[T comparable, R any](items []T, fn func(item T) (R, error)) (map[T]R, error) {
	ch := make(chan lo.Tuple2[T, R], len(items))
	ech := make(chan error, len(items))
	wg := sync.WaitGroup{}
	for _, item := range items {
		wg.Add(1)
		go func(arg T) {
			defer wg.Done()
			r, err := fn(arg)
			if err != nil {
				ech <- err
			} else {
				ch <- lo.T2(arg, r)
			}
		}(item)
	}
	wg.Wait()
	close(ch)
	close(ech)
	errs := lo.ChannelToSlice(ech)
	err := multierror.Append(nil, errs...).ErrorOrNil()
	r := lo.SliceToMap(lo.ChannelToSlice(ch), func(t lo.Tuple2[T, R]) (T, R) {
		return t.A, t.B
	})
	return r, err
}


