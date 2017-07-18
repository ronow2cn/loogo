package main

import (
	"client/app"
	"client/handler"
	"comm"
	"comm/config"
	"comm/logger"
	"comm/sched/loop"
	"flag"
	"os"
)

var log = logger.DefaultLogger

func main() {
	// load config
	config.Parse("config.json", "gate1")
	// open log
	logger.Open("")
	// signal
	quit := make(chan int)
	comm.OnSignal(func(s os.Signal) {
		log.Warning("shutdown signal received ...")
		close(quit)
	})
	// start
	N := flag.Int("n", 1, "client count")
	flag.Parse()

	handler.Init()
	app.ClientMgr.Start(*N)

	loop.Run()

	<-quit
	// stop
	app.ClientMgr.Stop()

	loop.Stop()
	// close log
	log.Notice("client stopped")
	logger.Close()
}
