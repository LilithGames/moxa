package service

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
)

var RouteNotFound codes.Code = 10000404

const DragonboatBalancerRouteShard = "DragonboatBalancerRouteShard"
const DragonboatBalancerRouteShardLeader = "DragonboatBalancerRouteShardLeader"
const DragonboatBalancerRouteNodeHost = "DragonboatBalancerRouteNodeHost"

func WithRouteShard(ctx context.Context, shardID uint64) context.Context {
	return context.WithValue(ctx, DragonboatBalancerRouteShard, shardID)
}
func WithRouteShardLeader(ctx context.Context, shardID uint64) context.Context {
	return context.WithValue(ctx, DragonboatBalancerRouteShardLeader, shardID)
}
func WithRouteNodeHost(ctx context.Context, nodeHostID string) context.Context {
	return context.WithValue(ctx, DragonboatBalancerRouteNodeHost, nodeHostID)
}

func init() {
	balancer.Register(NewDragonboatRoundRobinBalancerBuilder())
}

func NewDragonboatRoundRobinBalancerBuilder() balancer.Builder {
	builder := &dragonboatRoundRobinBalancerBuilder{}
	return base.NewBalancerBuilder("DragonboatRoundRobinBalancer", builder, base.Config{HealthCheck: true})
}

type dragonboatRoundRobinBalancerBuilder struct{}

func (it *dragonboatRoundRobinBalancerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	shards := make(map[uint64]*base.PickerBuildInfo, 0)
	leaders := make(map[uint64]*base.PickerBuildInfo, 0)
	nodeHosts := make(map[string]*base.PickerBuildInfo, 0)
	// log.Println("[INFO]", fmt.Sprintf("balancer start"))
	for conn, info := range buildInfo.ReadySCs {
		meta, ok := info.Address.Attributes.Value("meta").(*ResolverMeta)
		if !ok {
			log.Println("[WARN] dragonboatRoundRobinBalancerBuilder missing meta in Attributes, pls use with dragonboat resolver")
			continue
		}
		// log.Println("[INFO]", fmt.Sprintf("balancer: %s %v", info.Address.Addr, meta.ShardIDs))
		nodeHosts[info.Address.ServerName] = &base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{conn: info}}
		for shardID, isLeader := range meta.ShardIDs {
			if pbi, ok := shards[shardID]; ok {
				pbi.ReadySCs[conn] = info
			} else {
				shards[shardID] = &base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{conn: info}}
			}
			if isLeader {
				leaders[shardID] = &base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{conn: info}}
			}
		}
	}
	// log.Println("[INFO]", fmt.Sprintf("balancer end"))

	rrBuilder := &rrPickerBuilder{}
	sBuilder := &sPickerBuilder{}
	shardPickers := make(map[uint64]balancer.Picker, 0)
	leaderPickers := make(map[uint64]balancer.Picker, 0)
	nodeHostsPickers := make(map[string]balancer.Picker, 0)
	for shardID, pbi := range shards {
		shardPickers[shardID] = rrBuilder.Build(*pbi)
	}
	for shardID, pbi := range leaders {
		leaderPickers[shardID] = sBuilder.Build(*pbi)
	}
	for nodeHostID, pbi := range nodeHosts {
		nodeHostsPickers[nodeHostID] = sBuilder.Build(*pbi)
	}
	fallback := rrBuilder.Build(buildInfo)

	picker := &dragonboatRoundRobinBalancerPicker{
		shardPickers:     shardPickers,
		leaderPickers:    leaderPickers,
		nodeHostsPickers: nodeHostsPickers,
		fallback:         fallback,
	}
	return picker
}

type dragonboatRoundRobinBalancerPicker struct {
	shardPickers     map[uint64]balancer.Picker
	leaderPickers    map[uint64]balancer.Picker
	nodeHostsPickers map[string]balancer.Picker
	fallback         balancer.Picker
}

func (it *dragonboatRoundRobinBalancerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if nhid, ok := info.Ctx.Value(DragonboatBalancerRouteNodeHost).(string); ok {
		if picker, ok := it.nodeHostsPickers[nhid]; ok {
			return picker.Pick(info)
		} else {
			return balancer.PickResult{}, grpc.Errorf(RouteNotFound, "balancer nodehost %s not found", nhid)
		}
	} else if shardID, ok := info.Ctx.Value(DragonboatBalancerRouteShardLeader).(uint64); ok {
		if picker, ok := it.leaderPickers[shardID]; ok {
			return picker.Pick(info)
		} else {
			return balancer.PickResult{}, grpc.Errorf(RouteNotFound, "balancer shard leader %d not found", shardID)
		}
	} else if shardID, ok := info.Ctx.Value(DragonboatBalancerRouteShard).(uint64); ok {
		if picker, ok := it.shardPickers[shardID]; ok {
			return picker.Pick(info)
		} else {
			return balancer.PickResult{}, grpc.Errorf(RouteNotFound, "balancer shard %d not found", shardID)
		}
	} else {
		return it.fallback.Pick(info)
	}
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &rrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: rand.Intn(len(scs)),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}

type sPickerBuilder struct{}

func (*sPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	for sc := range info.ReadySCs {
		return &sPicker{conn: sc}
	}
	return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
}

type sPicker struct {
	conn balancer.SubConn
}

func (it *sPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{SubConn: it.conn}, nil
}
