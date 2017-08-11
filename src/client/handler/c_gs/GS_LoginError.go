package c_gs

import (
	"client/msg"
)

func GS_LoginError(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_LoginError)

	log.Info("GS_LoginError:", req.ErrorCode)
}
