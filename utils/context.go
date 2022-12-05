package utils

import (
	"context"

	"github.com/lni/goutils/syncutil"
)

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
