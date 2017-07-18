package tcp

import (
	"net"
	"time"
)

// ============================================================================

func Connect(addr string, timeout int32, f func(err error, sock *Socket)) {
	go func() {
		c, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			if f != nil {
				f(err, nil)
			}
			return
		}

		// create socket
		sock := newSocket(c)

		// parse remote addr
		sock.parseRemoteAddr()

		// event: connect
		if f != nil {
			f(nil, sock)
		}

		// go rw threads
		go sock.thrRead()
		go sock.thrWrite()
	}()
}
