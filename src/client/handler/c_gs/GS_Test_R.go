package c_gs

import (
	"client/msg"
)

func GS_Test_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Test_R)

	log.Info("Test Res:", req.Result)
}
