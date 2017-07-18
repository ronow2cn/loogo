package msg

// msgid ragne for C <-> GS: [100, 200)

// 测试
type C_Test struct { // msgid: 100
	Value int32
}

type GS_Test_R struct { // msgid: 101
	Result int32
}
