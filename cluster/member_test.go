package cluster

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	assert "github.com/stretchr/testify/require"
)

func TestMemberlist(t *testing.T) {
	start := time.Now()
	conf := memberlist.DefaultLocalConfig()
	conf.BindAddr = "127.0.0.1"
	ml, err := memberlist.Create(conf)
	assert.Nil(t, err)
	_ = ml
	fmt.Printf("duration: %s\n", time.Since(start))
}
