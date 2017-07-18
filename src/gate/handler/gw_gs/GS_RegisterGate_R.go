package gw_gs

import (
	"gate/app"
	"gate/msg"
)

func GS_RegisterGate_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_RegisterGate_R)
	gs := ctx.(*app.SocketGS)

	if req.Success {
		log.Notice("register to gs OK")
	} else {
		log.Notice("register to gs Failed")
		gs.Close()
	}
}
