package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"syscall"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/lni/goutils/syncutil"

	// "github.com/lni/dragonboat/v3/statemachine"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
	"github.com/lni/dragonboat/v3/config"
	"google.golang.org/grpc"

	"github.com/LilithGames/moxa"
	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/master_shard"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/sub_service"
	"github.com/LilithGames/moxa/sub_shard"
	"github.com/LilithGames/moxa/utils"
)

var masterShardVersion = ""
var subShardVersion = ""

func main() {
	stopper := syncutil.NewStopper()
	ctx := utils.BindContext(stopper, context.Background())
	utils.SignalHandler(stopper, syscall.SIGINT, syscall.SIGTERM)

	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	})
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	conf := &cluster.Config{
		NodeHostDir:     "/data/nodehost",
		LocalStorageDir: "/data/storage",
		MemberSeed:      []string{"moxa-headless:7946"},
		RttMillisecond:  100,
		DeploymentId:    0,
		EnableMetrics:   false,
	}
	sconf := &service.Config{
		HttpPort: 8000,
		GrpcPort: 8001,
	}
	master := &moxa.ShardProfile{
		Name:     moxa.MasterProfileName,
		CreateFn: runtime.NewMigrationStateMachineWrapper(master_shard.NewStateMachine),
		Config: config.Config{
			CheckQuorum:         true,
			ElectionRTT:         10,
			HeartbeatRTT:        1,
			SnapshotEntries:     100,
			CompactionOverhead:  0,
			OrderedConfigChange: true,
		},
		Version: masterShardVersion,
	}
	subProfileName := "SubProfile"
	sub := &moxa.ShardProfile{
		Name:     subProfileName,
		CreateFn: runtime.NewMigrationStateMachineWrapper(sub_shard.NewExampleStateMachine),
		Config: config.Config{
			CheckQuorum:         true,
			ElectionRTT:         10,
			HeartbeatRTT:        1,
			SnapshotEntries:     100,
			CompactionOverhead:  0,
			OrderedConfigChange: true,
		},
		Version: subShardVersion,
	}
	log.Println("[INFO]", fmt.Sprintf("MasterShardVersion: %s", masterShardVersion))
	log.Println("[INFO]", fmt.Sprintf("SubShardVersion: %s", subShardVersion))

	cm, err := cluster.NewKubernetesManager(ctx, conf)
	if err != nil {
		log.Fatalln("NewKubernetesClusterManager err: ", err)
	}
	sc, err := service.NewClient(cm, fmt.Sprintf("dragonboat://:%d", sconf.GrpcPort))
	if err != nil {
		log.Fatalln("NewServiceClient err: ", err)
	}
	if err := sc.Wait(ctx, time.Second*3); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("[WARN]", fmt.Errorf("ServiceClient.Wait err: %w", err))
		} else {
			log.Fatalln("[FATAL]", fmt.Errorf("ServiceClient.Wait err: %w", err))
		}
	}
	svc, err := service.NewClusterService(cm, sconf, service.ClusterServiceGrpcOpts(grpc.UnaryInterceptor(service.NewRequestTimeout(time.Second*60))))
	if err != nil {
		log.Fatalln("NewClusterService err: ", err)
	}
	service.RegisterNodeHostService(svc)
	service.RegisterSpecService(svc, sc)
	sub_service.RegisterApiService(svc)
	stopper.RunWorker(func() {
		if err := svc.Run(); err != nil {
			log.Println("[ERROR]", fmt.Errorf("ClusterService.Run err: %w", err))
			stopper.Close()
		}
	})

	sm, err := moxa.NewShardManager(cm, sc, master, sub)
	if err != nil {
		log.Fatalln("NewShardManager err: ", err)
	}
	go http.ListenAndServe(":6060", nil)
	stopper.RunWorker(func() {
		if err := cluster.StartWorker(stopper, cluster.UpdateNodeHostInfoWorker(cm)); err != nil {
			log.Println("[ERROR]", fmt.Errorf("StartClusterWorker err: %w", err))
			stopper.Close()
			return
		}
		if err := sm.RecoveryShards(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("[ERROR]", fmt.Errorf("ShardManager.RecoveryShards err: %w", err))
			stopper.Close()
			return
		}
		if err := sm.WaitServeMasterShard(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("[ERROR]", fmt.Errorf("ShardManager.ServeMasterShard err: %w", err))
			stopper.Close()
			return
		}
		if err := cluster.StartWorker(stopper,
			cluster.ShardSpecUpdatingWorker(cm),
			cluster.ShardSpecChangingWorker(cm),
			cluster.MigrationChangedWorker(cm),
			moxa.ShardSpecMembershipWorker(cm, sm),
			moxa.ShardSpecFixedSizeCreationWorker(cm, subProfileName, 1),
			moxa.MigrationNodeWorker(cm, sm),
		); err != nil {
			log.Println("[ERROR]", fmt.Errorf("StartClusterWorker err: %w", err))
			stopper.Close()
			return
		}
		cm.StartupReady().Set()
	})

	stopper.RunWorker(func() {
		<-stopper.ShouldStop()
		// drain
		drainCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := sm.Drain(drainCtx); err != nil {
			log.Println("[WARN]", fmt.Errorf("ShardManager.Drain() err: %w", err))
		}
		time.Sleep(time.Second)

		// terminate
		sm.Stop()
		svc.Stop()
		cm.Stop()
		log.Println("[INFO]", fmt.Sprintf("all system fully stopped"))
	})
	stopper.Wait()
	log.Println("[INFO]", fmt.Sprintf("system stopped as expected"))
}
