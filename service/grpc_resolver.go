package service

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
	"context"

	"github.com/lni/goutils/syncutil"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"github.com/hashicorp/go-multierror"

	"github.com/LilithGames/moxa/cluster"
)

type IMemberStateReader interface {
	GetMemberStateList() (map[string]*cluster.MemberState, uint64)
	GetMemberStateVersion() uint64
}

type serviceResolverBuilder struct {
	reader IMemberStateReader
}

func NewServiceResolverBuilder(reader IMemberStateReader) resolver.Builder {
	return &serviceResolverBuilder{reader: reader}
}

func (it *serviceResolverBuilder) Scheme() string {
	return "dragonboat"
}

func (it *serviceResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.URL.Port() == "" {
		return nil, fmt.Errorf("invalid resolver target(%s): port required", target.URL.String())
	}
	resolver := NewServiceResolver(target, cc, it.reader)
	return resolver, nil
}

type serviceResolver struct {
	conn    resolver.ClientConn
	reader IMemberStateReader
	version uint64
	target  resolver.Target
	stopper *syncutil.Stopper
	trigger chan struct{}
}

func NewServiceResolver(target resolver.Target, conn resolver.ClientConn, reader IMemberStateReader) *serviceResolver {
	it := &serviceResolver{
		conn:    conn,
		target:  target,
		reader:  reader,
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

func makeResolveMeta(m *cluster.MemberState) *ResolverMeta {
	shardIDs := map[uint64]bool{}
	for shardID, shard := range m.Shards {
		shardIDs[shardID] = shard.IsLeader
	}
	return &ResolverMeta{ShardIDs: shardIDs}
}

func resolveRaftAddress(raftAddress string) (string, error) {
	host, _, err := net.SplitHostPort(raftAddress)
	if err != nil {
		return "", fmt.Errorf("net.SplitHostPort(%s) err: %w", raftAddress, err)
	}
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", host)
	if err != nil {
		return "", fmt.Errorf("net.DefaultResolver.LookupIP(%s) err: %w", host, err)
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("resolve result empty")
	}
	ip := ips[0]
	return ip.String(), nil
}

func (it *serviceResolver) resolve() error {
	var errs error
	if it.version == 0 || it.version != it.reader.GetMemberStateVersion() {
		memberState, version := it.reader.GetMemberStateList()
		addrs := make([]resolver.Address, 0, len(memberState))
		for _, m := range memberState {
			ip, err := resolveRaftAddress(m.Meta.RaftAddress)
			if err != nil {
				errs = multierror.Append(errs, fmt.Errorf("resolveRaftAddress err: %w", err))
				continue
			}
			addr := resolver.Address{
				Addr:       fmt.Sprintf("%s:%s", ip, it.target.URL.Port()),
				ServerName: m.NodeHostId,
				Attributes: attributes.New("meta", makeResolveMeta(m)),
			}
			addrs = append(addrs, addr)
			it.conn.UpdateState(resolver.State{Addresses: addrs})
		}
		if errs == nil {
			it.version = version
		}
	}
	return errs
}

func (it *serviceResolver) watcher() {
	interval := time.Second * 3
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for {
		err := it.resolve()
		if err != nil {
			log.Println("[WARN]", fmt.Errorf("serviceResolver.watcher resolve err: %w", err))
			interval = time.Second
		} else {
			interval = time.Second * 3
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
