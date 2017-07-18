package tcp

import (
	"net"
	"sync"
)

// ============================================================================

type Server struct {
	lsn          net.Listener       // underlying listener
	cbConnection func(sock *Socket) // callback: connection
	cbError      func(err error)    // callback: error

	wg sync.WaitGroup
}

// ============================================================================

func CreateServer() *Server {
	return &Server{}
}

// ============================================================================

func (self *Server) Listen(addr string) *Server {
	if self.lsn != nil {
		return self
	}

	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		if self.cbError != nil {
			self.cbError(err)
		}
		return self
	}

	self.lsn = lsn

	self.wg.Add(1)
	go self.thrAccept()

	return self
}

func (self *Server) Stop() {
	if self.lsn == nil {
		return
	}

	self.lsn.Close()
	self.wg.Wait()
}

func (self *Server) OnConnection(f func(sock *Socket)) *Server {
	self.cbConnection = f
	return self
}

func (self *Server) OnError(f func(err error)) *Server {
	self.cbError = f
	return self
}

// ============================================================================

func (self *Server) thrAccept() {
	defer func() {
		self.lsn.Close()
		self.lsn = nil
		self.wg.Done()
	}()

	for {
		// accept
		c, err := self.lsn.Accept()
		if err != nil {
			break
		}

		// create socket
		sock := newSocket(c)

		// parse remote addr
		sock.parseRemoteAddr()

		// event: connection
		if self.cbConnection != nil {
			self.cbConnection(sock)
		}

		// go rw threads
		go sock.thrRead()
		go sock.thrWrite()
	}
}
