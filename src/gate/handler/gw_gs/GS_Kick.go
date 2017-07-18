package gw_gs

import (
	"gate/app"
	"gate/msg"
)

func GS_Kick(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Kick)

	app.NetMgr.KickSession(req.Sid)
}
