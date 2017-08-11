package c_gs

import (
	"game/app"
	"game/msg"
)

func C_Test(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_Test)
	plr := ctx.(*app.Player)

	plr.SendMsg(&msg.GS_Test_R{
		Result: req.Value + 1,
	})
}
