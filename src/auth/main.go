package main

import (
	"auth/app/dbmgr"
	"auth/route"
	"comm"
	"comm/config"
	"comm/logger"
	"comm/sched/asyncop"
	"comm/sched/loop"
	"flag"
	"os"
)

var log = logger.DefaultLogger

func main() {
	// parse command line
	argFile := flag.String("config", "config.json", "config file")
	argServer := flag.String("server", "auth", "server name")
	argLog := flag.String("log", "auth.log", "log file")

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
	// open db mgr
	dbmgr.Open()
	//loop
	loop.Run()
	asyncop.Start()

	//route
	route.AuthRoutes()

	// app started
	log.Notice("auth svr started")
}

func stop() {
	// close db mgr
	dbmgr.Close()
	// app stopped
	log.Notice("auth svr stopped")
}
