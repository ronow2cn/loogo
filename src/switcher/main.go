package main

import (
	"comm"
	"comm/config"
	"comm/logger"
	"comm/sched/loop"
	"flag"
	"os"
	"switcher/app"
	"switcher/app/dbmgr"
	"switcher/app/route"
)

var log = logger.DefaultLogger

func main() {
	// parse command line
	argFile := flag.String("config", "config.json", "config file")
	argServer := flag.String("server", "switcher", "server name")
	argLog := flag.String("log", "switcher.log", "log file")

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
	//routes
	route.Routes()
	//gate timer
	app.GatesTimer.Init()
	// app started
	log.Notice("switcher started")
}

func stop() {
	// close db mgr
	dbmgr.Close()
	// app stopped
	log.Notice("switcher stopped")
}
