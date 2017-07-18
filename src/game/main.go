package main

import (
	"comm"
	"comm/config"
	"comm/logger"
	"comm/sched/loop"
	"flag"
	"game/app"
	"game/handler"
	"game/perfmon"
	"math/rand"
	"os"
	"time"
)

var log = logger.DefaultLogger

func main() {
	rand.Seed(time.Now().Unix())
	// parse command line
	argFile := flag.String("config", "config.json", "config file")
	argServer := flag.String("server", "game1", "server name")
	argLog := flag.String("log", "game1.log", "log file")

	flag.Parse()
	// load config
	config.Parse(*argFile, *argServer)
	// open log
	logger.Open(*argLog)

	// signal
	quit := make(chan int)
	comm.OnSignal(func(s os.Signal) {
		log.Warning("shutdown signal received ...")
		close(quit)
	})

	start()
	<-quit
	stop()

	// close log
	logger.Close()
}

func start() {

	// start net mgr
	handler.Init()
	app.NetMgr.Start()
	// run app loop
	loop.Run()
	// app started
	log.Notice("game started")
	// perfmon
	perfmon.Start()
}

func stop() {
	// stop net mgr
	app.NetMgr.Stop()

	// stop app loop
	loop.Stop()

	// app stopped
	log.Notice("game stopped")
}
