package app

import (
	"comm/packet"
	"comm/sched/loop"
	"comm/tcp"
	"game/msg"
	"sync/atomic"
)

// ============================================================================

const (
	CMAXGATEID = 100000
)

var (
	COPREGISTERGATE = (&msg.GW_RegisterGate{}).MsgId()
)
var seqGateid int32 = CMAXGATEID
var TempSeqid = uint64(1)

// ============================================================================

type SocketGW struct {
	id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func newSocketGW(sock *tcp.Socket) *SocketGW {
	return &SocketGW{
		id:      atomic.AddInt32(&seqGateid, 1),
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketGW) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketGW) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(message.MsgId(), body)
	p.AddSid(0)

	self.SendPacket(p)
}

func (self *SocketGW) Close() {
	self.sock.Close()
}

//recieve gate msg, if sid > 0 mean this msg from client. sid == 0 mean just from gate
func (self *SocketGW) Dispatch(p packet.Packet) {
	// !Note: in net-thread
	op := p.Op()

	// check
	if !self.IsRegistered() && op != COPREGISTERGATE {
		return
	}

	// get msg creator
	f := msg.MsgCreators[op]
	if f == nil {
		return
	}

	// remove sid
	sid := p.RemoveSid()
	TempSeqid = sid
	// unmarshal
	message := f()
	err := msg.Unmarshal(p.Body(), message)
	if err != nil {
		return
	}

	// find handler
	h := msg.MsgHandlers[op]
	if h != nil {
		loop.Push(func() {
			// set ctx
			var ctx interface{}
			if sid == 0 {
				// directly from gate
				ctx = self
			} else {
				// from player
				ctx = PlayerMgr.FindPlayerBySid(sid) // this function MUST be run in loop thread
				if ctx == nil {
					return
				}
			}

			h(message, ctx)
		})
	}
}

func (self *SocketGW) IsRegistered() bool {
	return self.id < CMAXGATEID
}

func (self *SocketGW) GetTempId() uint64 {
	return TempSeqid
}
