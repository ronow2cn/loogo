package app

import (
	"comm/config"
	"comm/packet"
	"comm/tcp"
	"gate/msg"
)

// ============================================================================
//to gamesvr
type SocketGS struct {
	id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func newSocketGS(sock *tcp.Socket) *SocketGS {
	return &SocketGS{
		id:      config.DefaultGate.Id, //gate id
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketGS) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketGS) SendMsg(message msg.Message) {

	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	// assemble
	p := packet.Assemble(message.MsgId(), body)
	p.AddSid(0)

	// send
	self.SendPacket(p)
}

func (self *SocketGS) Close() {
	self.sock.Close()
}

//recieve msg from game, if sid > 0 mean this msg need send to client, sid == 0 mean not send to client
func (self *SocketGS) Dispatch(p packet.Packet) {
	// !Note: in net-thread
	sid := p.RemoveSid()
	if sid == 0 {
		// gate local msg
		op := p.Op()
		f := msg.MsgCreators[op]
		if f == nil {
			return
		}

		// unmarshal
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			log.Error("unmarshal msg failed:", err)
			self.Close()
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			h(message, self)
		}
	} else {
		// forward to session
		NetMgr.Forward2Session(sid, p)
	}
}
