package app

import (
	"comm/config"
	"comm/packet"
	"comm/tcp"
	"gate/msg"
	"sync"
	"sync/atomic"
)

// ============================================================================

var seqSid uint64

func InitSession() {
	seqSid = uint64(config.DefaultGate.Id) << 42
}

// ============================================================================

const (
	SessionStateNone = iota
	SessionStateAuthenticating
	SessionStateLoggedIn
	SessionStateLoggedOut
)

//come from client
type Session struct {
	id      uint64
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer

	locker sync.Mutex
	state  int
}

// ============================================================================

func newSession(sock *tcp.Socket) *Session {
	return &Session{
		id:      atomic.AddUint64(&seqSid, 1),
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
		state:   SessionStateNone,
	}
}

func (self *Session) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

//send msg to client
func (self *Session) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(message.MsgId(), body)

	self.SendPacket(p)
}

func (self *Session) Close() {
	self.sock.Close()
}

//recieve client msg, dispatch
func (self *Session) Dispatch(p packet.Packet) {
	op := p.Op()
	f := msg.MsgCreators[op]
	if f == nil {
		// op NOT found. forward to gs
		if self.state != SessionStateLoggedIn { // no need to lock
			return
		}

		NetMgr.Forward2GS(self.id, p)

	} else {
		// gate local handler
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			self.Close()
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			h(message, self)
		}
	}
}

// ============================================================================

func (self *Session) GetId() uint64 {
	return self.id
}

func (self *Session) GetIP() string {
	return self.sock.RemoteIP()
}

func (self *Session) BeginAuth() bool {
	self.locker.Lock()
	defer self.locker.Unlock()

	if self.state != SessionStateNone {
		return false
	}

	self.state = SessionStateAuthenticating
	return true
}

func (self *Session) EndAuth() {
	self.locker.Lock()
	defer self.locker.Unlock()

	if self.state == SessionStateAuthenticating {
		self.state = SessionStateNone
	}
}

func (self *Session) LoginPlayer(m *msg.GW_UserOnline) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if self.state != SessionStateAuthenticating {
		return
	}

	NetMgr.Send2GS(m)
	self.state = SessionStateLoggedIn
}

func (self *Session) LogoutPlayer() {
	self.locker.Lock()
	defer self.locker.Unlock()

	if self.state == SessionStateLoggedIn {
		NetMgr.Send2GS(&msg.GW_LogoutPlayer{
			Sid: self.id,
		})
	}

	self.state = SessionStateLoggedOut
}
