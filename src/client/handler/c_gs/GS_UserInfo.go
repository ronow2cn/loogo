package c_gs

import (
	"client/msg"
)

func GS_UserInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_UserInfo)

	log.Info("GS_UserInfo:", req)
}
