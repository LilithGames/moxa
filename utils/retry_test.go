package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	// assert "github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	bo := &Backoff{
		InitialDelay: time.Millisecond,
		Step:         2,
		MaxTimes:     lo.ToPtr(10),
		MaxDelay:     lo.ToPtr(time.Millisecond*20),
		MaxDuration:  lo.ToPtr(time.Millisecond*100),
	}
	for {
		d := bo.Next()
		if d == nil {
			return
		}
		fmt.Printf("%s\n", *d)
	}
}
