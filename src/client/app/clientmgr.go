package app

import (
	"comm/config"
	"comm/packet"
	"comm/tcp"
	"fmt"
	"sync"
	"time"
)

// ============================================================================

var ClientMgr = &clientmgrT{
	clients: make(map[int32]*Client),
}

var ClientNum int = 1

// ============================================================================

type clientmgrT struct {
	clients map[int32]*Client

	locker sync.Mutex
	wg     sync.WaitGroup
}

// ============================================================================

func (self *clientmgrT) Start(cnt int) {
	addr := fmt.Sprintf("%s:%d", config.DefaultGate.IPWan, config.DefaultGate.Port)
	ClientNum = cnt

	for i := 0; i < cnt; i++ {
		i := i

		self.wg.Add(2)
		go func() {
			defer self.wg.Done()

			time.Sleep(time.Duration(i*100) * time.Millisecond)
			tcp.Connect(addr, 3000, func(err error, sock *tcp.Socket) {
				defer self.wg.Done()

				if err != nil {
					log.Error("connect to server failed:", err)
					return
				}

				// add session
				client := newClient(sock)
				self.addClient(client)

				// sock event: data
				sock.OnData(func(buf []byte) {
					var p packet.Packet
					var err error

					for len(buf) > 0 {
						// read packet
						p, buf, err = client.preader.Read(buf)
						if err != nil {
							log.Debug("reading packet failed:", err)
							sock.Close()
							return
						}

						// no packet yet
						if p == nil {
							return
						}

						// got packet. dispatch
						client.Dispatch(p)
					}
				})

				// sock event: close
				sock.OnClose(func() {
					// remove client
					self.removeClient(client)
				})
			})
		}()
	}
}

func (self *clientmgrT) Stop() {
	// wait for all connect i/o
	self.wg.Wait()

	// close all clients
	self.closeAllClients()

	// wait
	for {
		if self.ClientCount() == 0 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (self *clientmgrT) ClientCount() int {
	self.locker.Lock()
	defer self.locker.Unlock()

	return len(self.clients)
}

// ============================================================================

func (self *clientmgrT) addClient(client *Client) {
	self.locker.Lock()
	self.clients[client.Id] = client
	self.locker.Unlock()

	// client event: connected
	client.OnConnected()

}

func (self *clientmgrT) removeClient(client *Client) {
	self.locker.Lock()
	delete(self.clients, client.Id)
	self.locker.Unlock()

	// client event: disconnected
	client.OnDisconnected()

}

func (self *clientmgrT) closeAllClients() {
	self.locker.Lock()
	defer self.locker.Unlock()

	for _, client := range self.clients {
		client.Close()
	}
}
