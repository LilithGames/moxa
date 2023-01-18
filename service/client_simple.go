package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SimpleClient struct {
	conn *grpc.ClientConn
}

func NewSimpleClient(target string) (IClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial err: %w", err)
	}
	return &SimpleClient{conn}, nil
}

func (it *SimpleClient) NodeHost() NodeHostClient {
	return NewNodeHostClient(it.conn)
}

func (it *SimpleClient) Spec() SpecClient {
	return NewSpecClient(it.conn)
}

func (it *SimpleClient) Wait(ctx context.Context, timeout time.Duration) error {
	panic("NotImplemented")
}

func (it *SimpleClient) Close() error {
	return it.conn.Close()
}
