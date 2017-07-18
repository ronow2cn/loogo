package msg

var MsgHandlers = map[uint32]func(message Message, ctx interface{}){}

func Handler(msgid uint32, h func(message Message, ctx interface{})) {
    MsgHandlers[msgid] = h
}
