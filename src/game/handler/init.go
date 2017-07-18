package handler

import (
    "game/handler/c_gs"
    "game/handler/gw_gs"
    "game/msg"
)

func Init() {
    msg.Handler(2000, gw_gs.GW_RegisterGate)
    msg.Handler(2002, gw_gs.GW_UserOnline)
    msg.Handler(2003, gw_gs.GW_LogoutPlayer)
    msg.Handler(100, c_gs.C_Test)
}
