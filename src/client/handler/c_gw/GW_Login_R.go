package c_gw

import (
	"client/app"
	"client/msg"
	"proto/errorcode"
)

func GW_Login_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_Login_R)
	client := ctx.(*app.Client)
	if req.ErrorCode != Err.OK {
		log.Error("client auth failed:", client.Id, "ErrorCode:", req.ErrorCode)
		client.Close()
		return
	} else {
		client.SendMsg(&msg.C_Test{
			Value: 12,
		})
		log.Info("SendMsg C_Test:")
	}
	log.Info("GW_Login_R:", req.ErrorCode, req.AuthId)
}
