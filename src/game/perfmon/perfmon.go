package perfmon

import (
	"comm/config"
	"comm/logger"
	"comm/sched/loop"
	"game/app"
	"time"
)

// ============================================================================

var log = logger.DefaultLogger

// ============================================================================

const N = 10

// ============================================================================

func Start() {
	if config.Common.Perfmon != "true" {
		return
	}

	go func() {
		for {
			time.Sleep(N * time.Second)

			log.Infof("loaded: %5d, online: %5d, loopq: %6d, handle: %6d",
				app.PlayerMgr.NumLoaded(),
				app.PlayerMgr.NumOnline(),
				loop.QLen(),
				loop.NumHandled()/N,
			)
		}
	}()
}
