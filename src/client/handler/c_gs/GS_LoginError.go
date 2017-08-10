package c_gs

import (
    "client/msg"
)

func GS_LoginError(message msg.Message, ctx interface{}) {
    req := message.(*msg.GS_LoginError)
    req = req
}
