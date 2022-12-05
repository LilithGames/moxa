package cluster

import (
	"fmt"

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
