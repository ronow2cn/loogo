package msg

var MsgCreators = map[uint32]func() Message{
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
    3000: func() Message {
        return &GS_LoginError{}
    },
    3001: func() Message {
        return &GS_UserInfo{}
    },
    100: func() Message {
        return &C_Test{}
    },
    101: func() Message {
        return &GS_Test_R{}
    },
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

func (self *GS_LoginError) MsgId() uint32 {
    return 3000
}

func (self *GS_UserInfo) MsgId() uint32 {
    return 3001
}

func (self *C_Test) MsgId() uint32 {
    return 100
}

func (self *GS_Test_R) MsgId() uint32 {
    return 101
}
