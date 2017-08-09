package route

import (
	"comm"
	"comm/config"
	"fmt"
	"net/http"
	"proto/macrocode"
)

func AuthRoutes() {
	go func() {
		http.HandleFunc(fmt.Sprintf("/auth/%d", macrocode.ChannelType_WeiXin), HandlerWeiXinAuth)

		err := http.ListenAndServe(config.Auth.IP+":"+comm.I32toa(config.Auth.Port), nil)
		if err != nil {
			log.Error("auth service ListenAndServe:", err)
		}
	}()
}
