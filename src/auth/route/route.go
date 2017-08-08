package route

import (
	"comm"
	"comm/config"
	"net/http"
)

func AuthRoutes() {
	go func() {
		http.HandleFunc("/auth/weixinauth", HandlerWeixinAuth)

		err := http.ListenAndServe(config.Auth.IP+":"+comm.I32toa(config.Auth.Port), nil)
		if err != nil {
			log.Error("auth service ListenAndServe:", err)
		}
	}()
}
