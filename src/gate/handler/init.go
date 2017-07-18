package handler

import (
    "gate/handler/c_gw"
    "gate/handler/gw_gs"
    "gate/msg"
)

func Init() {
    msg.Handler(1000, c_gw.C_Login)
    msg.Handler(2001, gw_gs.GS_RegisterGate_R)
    msg.Handler(2004, gw_gs.GS_Kick)
}
