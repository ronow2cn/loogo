package app

import (
	"client/msg"
	"comm"
	"comm/config"
	"comm/logger"
	"comm/packet"
	"comm/sched/loop"
	"comm/tcp"
	"proto/macrocode"
	"sync/atomic"
)

// ============================================================================
var log = logger.DefaultLogger
var seqClientID int32 = 0

// ============================================================================

type Client struct {
	Id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func newClient(sock *tcp.Socket) *Client {
	return &Client{
		Id:      atomic.AddInt32(&seqClientID, 1),
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *Client) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *Client) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(message.MsgId(), body)

	self.SendPacket(p)
}

func (self *Client) Close() {
	self.sock.Close()
}

func (self *Client) Dispatch(p packet.Packet) {
	// !Note: in net-thread
	op := p.Op()
	f := msg.MsgCreators[op]
	if f != nil {
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			log.Warning("unmarshal msg failed:", err)
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			loop.Push(func() {
				h(message, self)
			})
		}
	}
}

// ============================================================================

func (self *Client) OnConnected() {
	log.Info("client connected:", self.Id)

	//connect gate and send msg to gate
	self.SendMsg(&msg.C_Login{
		AuthChannel: macrocode.ChannelType_WeiXin,
		AuthType:    macrocode.LoginType_WeiXinCode,
		AuthId:      comm.I32toa(self.Id),
		Svr0:        config.DefaultGame.Name,
		VerMajor:    config.Common.VerMajor,
		VerMinor:    config.Common.VerMinor,
		VerBuild:    config.Common.VerBuild,
	})
}
func (self *Client) OnDisconnected() {
	log.Info("client disconnected:", self.Id)
}
