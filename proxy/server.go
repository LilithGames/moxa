package proxy

import (
	"net"
	"fmt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/siderolabs/grpc-proxy/proxy"

	"github.com/LilithGames/moxa/service"
)

type Server struct {
	listenPort int
	target string
	gs *grpc.Server
	backends []proxy.Backend
	conn *grpc.ClientConn
	reader service.IMemberStateReader
}

func NewServer(listenPort int, target string, reader service.IMemberStateReader) *Server {
	it := &Server{listenPort: listenPort, target: target, reader: reader}
	it.gs = grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(it.director)),
	)
	conn, err := it.createConn(context.TODO())
	if err != nil {
		panic(err)
	}
	it.conn = conn
	backend := &proxy.SingleBackend{GetConn: it.getConn}
	it.backends = []proxy.Backend{backend}
	return it
}

func (it *Server) createConn(ctx context.Context) (*grpc.ClientConn, error) {
	builder := service.NewServiceResolverBuilder(it.reader)
	conn, err := grpc.DialContext(ctx,
		it.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithCodec(proxy.Codec()),
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"DragonboatRoundRobinBalancer"}`),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial err: %w", err)
	}
	return conn, nil
}

func (it *Server) getConn(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	outCtx := metadata.NewOutgoingContext(ctx, md.Copy())
	return outCtx, it.conn, nil
}

func (it *Server) director(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {
	return proxy.One2One, it.backends, nil
}

func (it *Server) GrpcServiceRegistrar() grpc.ServiceRegistrar {
	return it.gs
}

func (it *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", it.listenPort))
	if err != nil {
		return fmt.Errorf("net.Listen err: %w", err)
	}
	if err := it.gs.Serve(lis); err != nil {
		return fmt.Errorf("server.Serve err: %w", err)
	}
	return nil
}
