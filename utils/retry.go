package utils

import (
	"context"
	"time"
)

func RetryWithDelay(ctx context.Context, max int, delay time.Duration, fn func() (bool, error)) error {
	var err error
	var retry bool
	for i := 0; i < max; i++ {
		retry, err = fn()
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
