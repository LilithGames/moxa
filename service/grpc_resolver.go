package service

import (
	"fmt"
	"log"
	"time"
	"net"
	"reflect"

	"github.com/lni/goutils/syncutil"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/LilithGames/moxa/cluster"
)

type serviceResolverBuilder struct {
	cm cluster.Manager
}

func NewServiceResolverBuilder(cm cluster.Manager) resolver.Builder {
	return &serviceResolverBuilder{cm: cm}
}

func (it *serviceResolverBuilder) Scheme() string {
	return "dragonboat"
}

func (it *serviceResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.URL.Port() == "" {
		return nil, fmt.Errorf("invalid resolver target(%s): port required", target.URL.String())
	}
	resolver := NewServiceResolver(target, cc, it.cm)
	return resolver, nil
}

type serviceResolver struct {
	conn    resolver.ClientConn
	cm      cluster.Manager
	target  resolver.Target
	stopper *syncutil.Stopper
	trigger chan struct{}
}

func NewServiceResolver(target resolver.Target, conn resolver.ClientConn, cm cluster.Manager) *serviceResolver {
	it := &serviceResolver{
		conn:    conn,
		target:  target,
		cm:      cm,
		stopper: syncutil.NewStopper(),
		trigger: make(chan struct{}, 1),
	}
	it.stopper.RunWorker(it.watcher)
	return it
}

func (it *serviceResolver) getHealthEndpoints() (map[string]struct{}, error) {
	hostname := it.target.URL.Hostname()
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("net.LookupIP(%s) err: %w", hostname, err)
	}
	result := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		result[ip.String()] = struct{}{}
	}
	return result, nil
}

type ResolverMeta struct {
	ShardIDs map[uint64]bool
}

func (it *ResolverMeta) Equal(v interface{}) bool {
	if it == nil && v == nil {
		return true
	}
	o, ok := v.(*ResolverMeta)
	if !ok {
		return false
	}
	return reflect.DeepEqual(it.ShardIDs, o.ShardIDs)
}

func (it *serviceResolver) resolve() (*resolver.State, error) {
	addrs := make([]resolver.Address, 0, it.cm.Members().Nums())
	it.cm.Members().Foreach(func(m cluster.MemberNode) bool {
		shardIDs := map[uint64]bool{}
		if m.State != nil {
			for shardID, shard := range m.State.Shards {
				shardIDs[shardID] = shard.IsLeader
			}
		}
		addr := resolver.Address{
			Addr:               fmt.Sprintf("%s:%s", m.Node.Addr.String(), it.target.URL.Port()),
			ServerName:         m.Node.Name,
			Attributes:         attributes.New("meta", &ResolverMeta{ShardIDs: shardIDs}),
		}
		addrs = append(addrs, addr)
		return true
	})
	return &resolver.State{Addresses: addrs}, nil
}

func (it *serviceResolver) watcher() {
	interval := time.Second*3
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for {
		state, err := it.resolve()
		if err != nil {
			log.Println("[WARN]", fmt.Errorf("serviceResolver.watcher resolve err: %w", err))
			interval = time.Second
		} else {
			it.conn.UpdateState(*state)
			interval = time.Second*3
		}

		timer.Stop()
		select {
		case <-timer.C:
		default:
		}
		timer.Reset(interval)
		select {
		case <-it.trigger:
		case <-timer.C:
		case <-it.stopper.ShouldStop():
			return
		}
	}
}

func (it *serviceResolver) ResolveNow(o resolver.ResolveNowOptions) {
	select {
	case it.trigger <- struct{}{}:
	default:
	}
}

func (it *serviceResolver) Close() {
	it.stopper.Close()
}
