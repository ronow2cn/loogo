package c_gs

import (
	"comm/packet"
	"game/app"
	"game/msg"
)

func C_Test(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_Test)
	gw := ctx.(*app.SocketGW)

	messageres := &msg.GS_Test_R{
		Result: req.Value + 1,
	}

	body, err := msg.Marshal(messageres)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(messageres.MsgId(), body)
	p.AddSid(gw.GetTempId())

	gw.SendPacket(p)

	log.Info("C_Test:", req.Value)

}
