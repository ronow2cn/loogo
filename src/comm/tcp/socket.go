package tcp

import (
	"comm"
	"net"
	"time"
)

// ============================================================================

const (
	CREADBUFFSIZE = 4096
	CMAXWRITESIZE = 2048
)

// ============================================================================

type Socket struct {
	c       net.Conn         // underlying connection
	qw      chan []byte      // write queue
	cbData  func(buf []byte) // callback: data
	cbClose func()           // callback: close
	rip     string           // remote ip
	rport   int32            // remote port
	hb      uint32           // heart beat in milliseconds
}

// ============================================================================

func newSocket(c net.Conn) *Socket {
	return &Socket{
		c:  c,
		qw: make(chan []byte, 1024),
	}
}

// ============================================================================

func (self *Socket) Close() {
	self.c.Close()
}

func (self *Socket) Send(buf []byte) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	self.qw <- buf
}

func (self *Socket) TcpNoDelay(b bool) {
	if c, ok := self.c.(*net.TCPConn); ok {
		c.SetNoDelay(b)
	}
}

func (self *Socket) RemoteAddr() string {
	return self.c.RemoteAddr().String()
}

func (self *Socket) RemoteIP() string {
	return self.rip
}

func (self *Socket) RemotePort() int32 {
	return self.rport
}

func (self *Socket) HeartBeat(ms uint32) {
	self.hb = ms
}

func (self *Socket) OnData(f func(buf []byte)) {
	self.cbData = f
}

func (self *Socket) OnClose(f func()) {
	self.cbClose = f
}

// ============================================================================

func (self *Socket) parseRemoteAddr() {
	addr := self.c.RemoteAddr().String()
	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		self.rip = host
		self.rport = comm.Atoi32(port)
	}
}

func (self *Socket) thrRead() {
	defer func() {
		self.Close()
		close(self.qw)
		self.cbData = nil
		self.cbClose = nil
	}()

	c := self.c
	buf := make([]byte, CREADBUFFSIZE)

	for {
		// set deadline
		if self.hb > 0 {
			c.SetReadDeadline(time.Now().Add(time.Duration(self.hb) * time.Millisecond))
		}

		// read
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			// event: close
			if self.cbClose != nil {
				self.cbClose()
			}

			break
		}

		// event: data
		if self.cbData != nil {
			self.cbData(buf[:n])
		}
	}
}

func (self *Socket) thrWrite() {
	c := self.c
	qw := self.qw

	for {
		select {
		case buf, ok := <-qw:
			if !ok {
				return
			}

			L := len(qw)
			for L > 0 && len(buf) < CMAXWRITESIZE {
				buf = append(buf, <-qw...)
				L--
			}

			n, err := c.Write(buf)
			if err != nil {
				return
			}

			if n != len(buf) {
				panic("socket write error!")
			}
		}
	}
}
