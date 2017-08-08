package handler

import (
	"comm"
	"comm/config"
	"net/http"
)

func Routes() {

	go func() {
		http.HandleFunc("/server/gatereg", HandlerGateReg)
		http.HandleFunc("/server/gamelist", HandlerGameList)
		http.HandleFunc("/server/gatelist", HandlerGateList)
		http.HandleFunc("/server/usersvr", HandlerUserSvr)

		err := http.ListenAndServe(config.Switcher.IP+":"+comm.I32toa(config.Switcher.Port), nil)
		if err != nil {
			log.Error("switcher service ListenAndServe:", err)
		}
	}()
}
