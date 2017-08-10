package msg

// msgid ragne for C <-> GS: [3000, 3100)

// 登录错误
type GS_LoginError struct { // msgid: 3000
	ErrorCode int
}

// 玩家信息
type GS_UserInfo struct { // msgid: 3001
	UserId string
	Name   string
	Head   string
	Lv     int32
	Exp    int32
	Vip    int32
}
