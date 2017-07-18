package app

import (
	"comm"
	"comm/config"
	"comm/logger"
	"comm/packet"
	"comm/tcp"
	"fmt"
	"gate/msg"
	"sync"
	"time"
)

var log = logger.DefaultLogger

// ============================================================================

var NetMgr = &netmgrT{
	sessions: make(map[uint64]*Session),
	connectq: tcp.NewConnectQ(),
}

// ============================================================================

const (
	CMAXSESSIONCOUNT = 30000
)

// ============================================================================

type netmgrT struct {
	svr4c    *tcp.Server         // server for client
	sessions map[uint64]*Session // session map

	connectq *tcp.ConnectQ // connect queue
	cnnGS    *SocketGS     // connection to gs

	locker sync.Mutex
}

// ============================================================================

func (self *netmgrT) Start() {
	log.Info("starting net mgr ...")

	self.listenOnClient()

	self.connectq.Open()
	self.connectq.Connect(self.connectToGS, 0)
}

func (self *netmgrT) Stop() {
	log.Info("stopping net mgr ...")

	// stop servers
	self.svr4c.Stop()
	self.closeAllSessions()

	// stop connections
	self.connectq.Close()
	self.closeAllConnections()

	// wait
	for {
		if self.SessionCount() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

//msg: client >> gate >> gamesvr
func (self *netmgrT) Forward2GS(sid uint64, p packet.Packet) {
	cnnGS := self.cnnGS
	if cnnGS != nil {
		p.AddSid(sid)
		cnnGS.SendPacket(p)
	}
}

//msg: gamesvr >> gate >> client
func (self *netmgrT) Forward2Session(sid uint64, p packet.Packet) {
	sess := self.findSession(sid)
	if sess != nil {
		sess.SendPacket(p)
	}
}

//msg: gate >> gamesvr
func (self *netmgrT) Send2GS(message msg.Message) {
	cnnGS := self.cnnGS
	if cnnGS != nil {
		cnnGS.SendMsg(message)
	}
}

func (self *netmgrT) KickSession(sid uint64) {
	sess := self.findSession(sid)
	if sess != nil {
		sess.Close()
	}
}

func (self *netmgrT) SessionCount() int {
	self.locker.Lock()
	defer self.locker.Unlock()

	return len(self.sessions)
}

// ============================================================================

func (self *netmgrT) listenOnClient() {
	// some cloud-servers DO NOT allow us to listen on specific WAN IP
	// 	so, listen on 0.0.0.0
	addr := fmt.Sprintf("0.0.0.0:%d", config.DefaultGate.Port)

	self.svr4c = tcp.CreateServer().
		OnConnection(func(sock *tcp.Socket) {
			// check session count
			if self.SessionCount() > CMAXSESSIONCOUNT {
				sock.Close()
				return
			}

			// add session
			sess := newSession(sock)
			self.addSession(sess)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = sess.preader.Read(buf)
					if err != nil {
						log.Debug("reading packet failed:", sock.RemoteAddr(), err)
						sock.Close()
						return
					}

					// no packet yet
					if p == nil {
						return
					}

					// got packet. dispatch
					sess.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove session
				self.removeSession(sess)
			})
		}).
		OnError(func(err error) {
			comm.Panic("listen on client failed:", err)
		}).
		Listen(addr)

	log.Notice("listen on client:", addr)
}

func (self *netmgrT) connectToGS(done func()) {
	log.Info("connecting to gs ...")

	tcp.Connect(config.DefaultGame.Addr4GW, 3000, func(err error, sock *tcp.Socket) {
		defer done()

		if err != nil {
			log.Warning("connect to gs failed:", err)
			self.connectq.Connect(self.connectToGS, 5000)
			return
		}

		// create cnn gs
		cnnGS := newSocketGS(sock)
		self.cnnGS = cnnGS

		// sock event: data
		sock.OnData(func(buf []byte) {
			var p packet.Packet
			var err error

			for len(buf) > 0 {
				// read packet
				p, buf, err = cnnGS.preader.Read(buf)
				if err != nil {
					log.Debug("reading packet failed:", sock.RemoteAddr(), err)
					sock.Close()
					return
				}

				// no packet yet
				if p == nil {
					return
				}

				// got packet. dispatch
				cnnGS.Dispatch(p)
			}
		})

		// sock event: close
		sock.OnClose(func() {
			log.Warning("connection to gs disconnected")
			self.cnnGS = nil
			self.closeAllSessions()
			self.connectq.Connect(self.connectToGS, 5000)
			return
		})

		// register
		cnnGS.SendMsg(&msg.GW_RegisterGate{
			Id: cnnGS.id,
		})
	})
}

func (self *netmgrT) addSession(sess *Session) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.sessions[sess.id] = sess
}

func (self *netmgrT) removeSession(sess *Session) {
	self.locker.Lock()
	delete(self.sessions, sess.id)
	self.locker.Unlock()

	// logout player
	sess.LogoutPlayer()
}

func (self *netmgrT) findSession(sid uint64) *Session {
	self.locker.Lock()
	defer self.locker.Unlock()

	return self.sessions[sid]
}

func (self *netmgrT) closeAllSessions() {
	self.locker.Lock()
	defer self.locker.Unlock()

	for _, sess := range self.sessions {
		sess.Close()
	}
}

func (self *netmgrT) closeAllConnections() {
	// close cnnGS
	cnnGS := self.cnnGS
	if cnnGS != nil {
		cnnGS.Close()
	}
}
