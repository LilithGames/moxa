package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/lni/goutils/syncutil"
)

type Worker func(stopper *syncutil.Stopper) error

func StartWorker(stopper *syncutil.Stopper, workers ...Worker) error {
	for _, worker := range workers {
		if err := worker(stopper); err != nil {
			return fmt.Errorf("%v() err: %w", worker, err)
		}
	}
	return nil
}

func ResetTimer(t *time.Timer, d time.Duration) {
	t.Stop()
	select {
	case <-t.C:
	default:
	}
	t.Reset(d)
}

func BindContext(stopper *syncutil.Stopper, ctx context.Context) context.Context {
	cctx, cancel := context.WithCancel(ctx)
	stopper.RunWorker(func() {
		select {
		case <-stopper.ShouldStop():
			cancel()
		case <-cctx.Done():
			stopper.Close()
		}
	})
	return cctx
}
