package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/lni/goutils/syncutil"
	"github.com/lni/dragonboat/v3/statemachine"

	"github.com/LilithGames/moxa"
	"github.com/LilithGames/moxa/cluster"
	"github.com/LilithGames/moxa/service"
	"github.com/LilithGames/moxa/sub_shard"
	"github.com/LilithGames/moxa/sub_service"
)

func SignalHandler(stopper *syncutil.Stopper, signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	stopper.RunWorker(func() {
		select {
		case <-ch:
			stopper.Close()
		case <-stopper.ShouldStop():
			return
		}
	})
}

func main() {
	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	})

	cm, err := cluster.NewKubernetesManager()
	if err != nil {
		log.Fatalln("NewKubernetesClusterManager err: ", err)
	}
	sc, err := service.NewClient(cm)
	if err != nil {
		log.Fatalln("NewServiceClient err: ", err)
	}
	svc, err := service.NewClusterService(cm)
	if err != nil {
		log.Fatalln("NewClusterService err: ", err)
	}
	service.RegisterNodeHostService(svc)
	service.RegisterSpecService(svc)
	sub_service.RegisterApiService(svc)
	sm, err := moxa.NewShardManager(cm, sc, statemachine.CreateStateMachineFunc(sub_shard.NewExampleStateMachine))
	if err != nil {
		log.Fatalln("NewShardManager err: ", err)
	}

	stopper := syncutil.NewStopper()
	ctx := cluster.BindContext(stopper, context.Background())
	SignalHandler(stopper, syscall.SIGINT, syscall.SIGTERM)
	stopper.RunWorker(func() {
		if err := svc.Run(); err != nil {
			log.Println("[ERROR]", fmt.Errorf("ClusterService.Run err: %w", err))
			stopper.Close()
		}
	})
	stopper.RunWorker(func() {
		if err := sm.Migrate(ctx); err != nil {
			log.Println("[ERROR]", fmt.Errorf("ShardManager.Migrate err: %w", err))
			stopper.Close()
			return
		}
		if err := sm.RecoveryShards(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("[ERROR]", fmt.Errorf("ShardManager.RecoveryShards err: %w", err))
			stopper.Close()
			return
		}
		if err := cluster.StartWorker(stopper, cluster.UpdateNodeHostInfoWorker(cm)); err != nil {
			log.Println("[ERROR]", fmt.Errorf("StartClusterWorker err: %w", err))
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
			moxa.ShardSpecMembershipWorker(cm, sm),
		); err != nil {
			log.Println("[ERROR]", fmt.Errorf("StartClusterWorker err: %w", err))
			stopper.Close()
			return
		}
	})
	stopper.RunWorker(func() {
		<-stopper.ShouldStop()
		// drain
		sm.Drain()
		time.Sleep(time.Second)

		// terminate
		svc.Stop()
		cm.Stop()
	})
	stopper.Wait()
}
