package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"

	"github.com/LilithGames/moxa/service"
)

type IHelper interface {
	MustClient() service.IClient
	Ctx() context.Context
	Cancel() 
}

type Helper struct {
	cctx *cli.Context
	ctx  context.Context
	cancel context.CancelFunc
}

func NewHelper(cCtx *cli.Context) IHelper {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Helper{cCtx, ctx, cancel}
}

func (it *Helper) MustClient() service.IClient {
	target := it.cctx.String("address")
	client, err := service.NewSimpleClient(target)
	if err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("service.NewSimpleClient err: %w", err))
	}
	return client
}

func (it *Helper) Ctx() context.Context {
	return it.ctx
}

func (it *Helper) Cancel() {
	it.cancel()
}
