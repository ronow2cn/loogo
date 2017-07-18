package msg

var MsgCreators = map[uint32]func() Message{
    1000: func() Message {
        return &C_Login{}
    },
    1001: func() Message {
        return &GW_Login_R{}
    },
    2000: func() Message {
        return &GW_RegisterGate{}
    },
    2001: func() Message {
        return &GS_RegisterGate_R{}
    },
    2002: func() Message {
        return &GW_UserOnline{}
    },
    2003: func() Message {
        return &GW_LogoutPlayer{}
    },
    2004: func() Message {
        return &GS_Kick{}
    },
}

func (self *C_Login) MsgId() uint32 {
    return 1000
}

func (self *GW_Login_R) MsgId() uint32 {
    return 1001
}

func (self *GW_RegisterGate) MsgId() uint32 {
    return 2000
}

func (self *GS_RegisterGate_R) MsgId() uint32 {
    return 2001
}

func (self *GW_UserOnline) MsgId() uint32 {
    return 2002
}

func (self *GW_LogoutPlayer) MsgId() uint32 {
    return 2003
}

func (self *GS_Kick) MsgId() uint32 {
    return 2004
}
