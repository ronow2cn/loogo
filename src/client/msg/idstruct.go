package msg

var MsgCreators = map[uint32]func() Message{
    1000: func() Message {
        return &C_Login{}
    },
    1001: func() Message {
        return &GW_Login_R{}
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

func (self *C_Login) MsgId() uint32 {
    return 1000
}

func (self *GW_Login_R) MsgId() uint32 {
    return 1001
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
