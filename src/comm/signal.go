package comm

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSignal(f func(os.Signal)) {
	go func() {
		c := make(chan os.Signal, 8)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		for s := range c {
			f(s)
			if s != syscall.SIGHUP {
				break
			}
		}
	}()
}
