package msg

// msgid ragne for GW <-> GS: [2000, 3000)

// 注册网关
type GW_RegisterGate struct { // msgid: 2000
	Id int32 // 网关 Id
}

type GS_RegisterGate_R struct { // msgid: 2001
	Success bool
}

// 玩家上线
type GW_UserOnline struct { // msgid: 2002
	Channel    int32
	ChannelUid string
	Sid        uint64
	UserId     string
	Svr0       string
	LoginIP    string
}

// 通知 game 玩家登出
type GW_LogoutPlayer struct { // msgid: 2003
	Sid uint64
}

// game 踢人
type GS_Kick struct { // msgid: 2004
	Sid uint64
}
