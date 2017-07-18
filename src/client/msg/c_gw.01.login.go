package msg

// msgid ragne for C <-> GW: [1000, 2000)

// 登录
type C_Login struct { // msgid: 1000
	AuthChannel int32  // 渠道
	AuthType    int32  // 认证类型
	AuthId      string // 认证 id
	AuthToken   string // 认证 token
	Svr0        string // 初始服名称
	Param1      string // 额外参数 1
	Param2      string // 额外参数 2
	Param3      string // 额外参数 3
	VerMajor    string // 主版本号
	VerMinor    string // 次版本号
	VerBuild    string // build 版本号
}

// 登录回复
type GW_Login_R struct { // msgid: 1001
	ErrorCode int
	AuthId    string
	Token     string `msgpack:",omitempty"`
	ExpireT   int64  `msgpack:",omitempty"`
}
