package utils

import (
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/lni/goutils/syncutil"
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
