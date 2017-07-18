package main

import (
	"comm"
	"comm/config"
	"comm/logger"
	"flag"
	"gate/app"
	"gate/app/dbmgr"
	"gate/handler"
	"os"
)

var log = logger.DefaultLogger

func main() {
	// parse command line
	argFile := flag.String("config", "config.json", "config file")
	argServer := flag.String("server", "gate1", "server name")
	argLog := flag.String("log", "gate1.log", "log file")

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
	// init session
	app.InitSession()
	// init handler func
	handler.Init()
	// listen client and connect gs
	app.NetMgr.Start()
	// gwreg
	app.GWReg.Run()
	// app started
	log.Notice("gate started")
}

func stop() {
	// stop net mgr
	app.NetMgr.Stop()
	// close db mgr
	dbmgr.Close()
	// app stopped
	log.Notice("gate stopped")
}
