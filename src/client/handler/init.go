package handler

import (
    "client/handler/c_gs"
    "client/handler/c_gw"
    "client/msg"
)

func Init() {
    msg.Handler(1001, c_gw.GW_Login_R)
    msg.Handler(3000, c_gs.GS_LoginError)
    msg.Handler(3001, c_gs.GS_UserInfo)
    msg.Handler(101, c_gs.GS_Test_R)
}
