package service

import (
	"context"
	"fmt"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/LilithGames/moxa/cluster"
)

type IClient interface {
	Close() error
	Wait(ctx context.Context, timeout time.Duration) error
	NodeHost() NodeHostClient
	Spec() SpecClient
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(cm cluster.Manager, target string) (IClient, error) {
	builder := NewServiceResolverBuilder(cm.Members())

	retry_opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponentialWithJitter(100*time.Millisecond, 0.01)),
		grpc_retry.WithCodes(RouteNotFound),
	}

	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"DragonboatRoundRobinBalancer"}`),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retry_opts...)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial err: %w", err)
	}
	return &Client{conn}, nil
}

func (it *Client) NodeHost() NodeHostClient {
	return NewNodeHostClientWrapper(NewNodeHostClient(it.conn))
}

func (it *Client) Spec() SpecClient {
	return NewSpecClientWrapper(NewSpecClient(it.conn))
}

func (it *Client) Wait(ctx context.Context, timeout time.Duration) error {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		it.conn.Connect()
		state := it.conn.GetState()
		if state == connectivity.Ready {
			return nil
		} else if state == connectivity.Shutdown {
			return fmt.Errorf("connectivity.Shutdown")
		}
		if !it.conn.WaitForStateChange(cctx, state) {
			return fmt.Errorf("WaitForStateChange %w", cctx.Err())
		}
	}
}

func (it *Client) Close() error {
	return it.conn.Close()
}
