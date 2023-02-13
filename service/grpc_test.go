package service

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGrpcMetadata(t *testing.T) {
	ctx := context.TODO()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md.Set("name", "hulucc")
	} else {
		md = metadata.Pairs("name", "hulucc")
	}
	ctx1 := metadata.NewOutgoingContext(ctx, md)
	fmt.Printf("%+v\n", ctx1)
}
