package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/go-multierror"
	"github.com/lni/goutils/syncutil"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/utils"
)

type GrpcGatewayRegisterFunc func(ctx context.Context, mux *runtime.ServeMux) error

type ClusterService struct {
	config    *Config
	cm        cluster.Manager
	server    *http.Server
	gs        *grpc.Server
	stopper   *syncutil.Stopper
	registers []GrpcGatewayRegisterFunc
}

func NewClusterService(cm cluster.Manager, config *Config) (*ClusterService, error) {
	it := &ClusterService{
		config:    config,
		cm:        cm,
		server:    &http.Server{Addr: fmt.Sprintf(":%d", config.HttpPort)},
		gs:        grpc.NewServer(grpc.UnaryInterceptor(NewRequestTimeout(time.Second * 60))),
		stopper:   syncutil.NewStopper(),
		registers: []GrpcGatewayRegisterFunc{},
	}
	return it, nil
}

func (it *ClusterService) AddGrpcGatewayRegister(fn GrpcGatewayRegisterFunc) {
	it.registers = append(it.registers, fn)
}

func (it *ClusterService) GrpcServiceRegistrar() grpc.ServiceRegistrar {
	return it.gs
}

func (it *ClusterService) ClusterManager() cluster.Manager {
	return it.cm
}

func (it *ClusterService) Config() *Config {
	return it.config
}

func (it *ClusterService) Stopper() *syncutil.Stopper {
	return it.stopper
}

func (it *ClusterService) Run() error {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{Indent: "  ", Multiline: true, EmitUnpopulated: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	}))
	it.server.Handler = mux

	ech := make(chan error, 10)
	it.stopper.RunWorker(func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", it.config.GrpcPort))
		if err != nil {
			ech <- fmt.Errorf("net.Listen err: %w", err)
			it.stopper.Close()
			return
		}
		if err := it.gs.Serve(lis); err != nil {
			ech <- fmt.Errorf("grpc.Serve err: %w", err)
			it.stopper.Close()
			return
		}
		log.Println("[INFO]", fmt.Sprintf("grpc.Serve stopped"))
	})
	ctx := utils.BindContext(it.stopper, context.Background())
	it.stopper.RunWorker(func() {
		for _, register := range it.registers {
			if err := register(ctx, mux); err != nil {
				ech <- fmt.Errorf("GrpcGatewayRegister %v err: %w", register, err)
				it.stopper.Close()
				return
			}
		}
	})
	it.stopper.RunWorker(func() {
		if err := it.server.ListenAndServe(); err != http.ErrServerClosed {
			ech <- fmt.Errorf("http.ListenAndServe err: %w", err)
			it.stopper.Close()
			return
		}
		log.Println("[INFO]", fmt.Sprintf("http.ListenAndServe stopped"))
	})
	it.stopper.RunWorker(func() {
		<-it.stopper.ShouldStop()
		it.gs.GracefulStop()
		it.server.Shutdown(context.Background())
	})
	it.stopper.Wait()
	close(ech)
	err := multierror.Append(nil, lo.ChannelToSlice(ech)...).ErrorOrNil()
	return err
}

func (it *ClusterService) Stop() error {
	it.stopper.Stop()
	return nil
}
