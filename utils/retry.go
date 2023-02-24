package utils

import (
	"context"
	"time"
	"fmt"
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

func RetryWithBackoff(ctx context.Context, bo *Backoff, fn func() (bool, error)) error {
	b := bo.New()
	for {
		retry, err := fn()
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		delay := b.Next()
		if delay == nil {
			return fmt.Errorf("Backoff Reached: %w", err)
		}
		select {
		case <-time.After(*delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type Backoff struct {
	InitialDelay time.Duration
	Step         float32
	MaxDelay     *time.Duration
	MaxTimes     *int
	MaxDuration  *time.Duration

	times int
	last  time.Duration
	sum   time.Duration
}

func (it *Backoff) New() *Backoff {
	return &Backoff{
		InitialDelay: it.InitialDelay,
		Step: it.Step,
		MaxDelay: it.MaxDelay,
		MaxTimes: it.MaxTimes,
		MaxDuration: it.MaxDuration,
	}
}

func (it *Backoff) Next() *time.Duration {
	if it.times == 0 {
		it.times++
		it.sum += it.InitialDelay
		it.last = it.InitialDelay
		return &it.InitialDelay
	}
	delay := it.last * time.Duration(it.Step)
	if it.MaxDelay != nil {
		if delay > *it.MaxDelay {
			delay = *it.MaxDelay
		}
	}
	if it.MaxTimes != nil {
		if it.times > *it.MaxTimes {
			return nil
		}
	}
	if it.MaxDuration != nil {
		if it.sum > *it.MaxDuration {
			return nil
		}
	}
	it.times++
	it.sum += delay
	it.last = delay
	return &delay
}
