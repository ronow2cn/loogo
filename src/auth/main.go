package main

import (
	"comm"
	"comm/config"
	"comm/logger"
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

	//loop
	loop.Run()

	// app started
	log.Notice("switcher started")
}

func stop() {

	// app stopped
	log.Notice("switcher stopped")
}
