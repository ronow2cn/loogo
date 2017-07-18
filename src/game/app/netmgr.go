package app

import (
	"comm"
	"comm/config"
	"comm/packet"
	"comm/sched/loop"
	"comm/tcp"
	"game/msg"
	"sync"
	"time"
)

// ============================================================================

var NetMgr = &netmgrT{
	gates: make(map[int32]*SocketGW),
}

// ============================================================================

type netmgrT struct {
	svr4gw *tcp.Server         // server for gate
	gates  map[int32]*SocketGW // gate map

	locker sync.Mutex
}

// ============================================================================

func (self *netmgrT) Start() {
	log.Info("starting net mgr ...")

	self.listenOnGates()
}

func (self *netmgrT) Stop() {
	log.Info("stopping net mgr ...")

	// stop servers
	self.svr4gw.Stop()
	self.closeAllGates()

	// wait
	for {
		if self.GateCount() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *netmgrT) RegisterGate(gw *SocketGW, id int32) bool {
	self.locker.Lock()
	defer self.locker.Unlock()

	// check reg id
	if id >= CMAXGATEID {
		return false
	}

	// old entry MUST exist
	// new entry MUST NOT exist
	if self.gates[gw.id] == nil || self.gates[id] != nil {
		return false
	}

	// remove old entry
	delete(self.gates, gw.id)

	// set registered id
	gw.id = id

	// add new entry
	self.gates[gw.id] = gw

	log.Notice("gate registered:", id)

	return true
}

func (self *netmgrT) Send2Gate(gateid int32, message msg.Message) {
	if gateid == 0 {
		// broadcast
		for _, gw := range self.arrayAllGates() {
			gw.SendMsg(message)
		}
	} else {
		// specific one
		gw := self.findGate(gateid)
		if gw != nil {
			gw.SendMsg(message)
		}
	}
}

func (self *netmgrT) Send2Player(sid uint64, message msg.Message) {
	// find gate
	gateid := sid2gateid(sid)
	gw := self.findGate(gateid)
	if gw == nil {
		return
	}

	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	// assemble
	p := packet.Assemble(message.MsgId(), body)
	p.AddSid(sid)

	// send
	gw.SendPacket(p)
}

func (self *netmgrT) GateCount() int {
	self.locker.Lock()
	defer self.locker.Unlock()

	return len(self.gates)
}

// ============================================================================

func (self *netmgrT) listenOnGates() {
	self.svr4gw = tcp.CreateServer().
		OnConnection(func(sock *tcp.Socket) {
			// add gate
			gw := newSocketGW(sock)
			self.addGate(gw)

			// sock event: data
			sock.OnData(func(buf []byte) {
				var p packet.Packet
				var err error

				for len(buf) > 0 {
					// read packet
					p, buf, err = gw.preader.Read(buf)
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
					gw.Dispatch(p)
				}
			})

			// sock event: close
			sock.OnClose(func() {
				// remove gate
				self.removeGate(gw)
			})
		}).
		OnError(func(err error) {
			comm.Panic("listen on gate failed:", err)
		}).
		Listen(config.DefaultGame.Addr4GW)

	log.Notice("listen on gates:", config.DefaultGame.Addr4GW)
}

func (self *netmgrT) addGate(gw *SocketGW) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.gates[gw.id] = gw
}

func (self *netmgrT) removeGate(gw *SocketGW) {
	self.locker.Lock()
	delete(self.gates, gw.id)
	self.locker.Unlock()

	if gw.IsRegistered() {
		log.Warning("gate dropped:", gw.id)

		loop.Push(func() {
			// offline players from this gate
			PlayerMgr.OfflineAllPlayers(gw.id)
		})
	}
}

func (self *netmgrT) findGate(id int32) *SocketGW {
	self.locker.Lock()
	defer self.locker.Unlock()

	gw := self.gates[id]
	if gw != nil && gw.IsRegistered() {
		return gw
	} else {
		return nil
	}
}

func (self *netmgrT) arrayAllGates() (ret []*SocketGW) {
	self.locker.Lock()
	defer self.locker.Unlock()

	for _, gw := range self.gates {
		if gw.IsRegistered() {
			ret = append(ret, gw)
		}
	}

	return
}

func (self *netmgrT) closeAllGates() {
	self.locker.Lock()
	defer self.locker.Unlock()

	for _, gw := range self.gates {
		gw.Close()
	}
}
