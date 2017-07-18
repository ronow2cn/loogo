package msg

import (
    "gopkg.in/vmihailenco/msgpack.v2"
)

type Message interface {
    MsgId() uint32
}

func Marshal(m Message) ([]byte, error) {
    return msgpack.Marshal(m)
}

func Unmarshal(b []byte, obj interface{}) error {
    return msgpack.Unmarshal(b, obj)
}
