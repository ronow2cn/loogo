package tcp

import (
	"sync"
	"time"
)

// ============================================================================

type ConnectQ struct {
	q    chan *connectF
	wg   sync.WaitGroup
	quit chan int
}

type connectF struct {
	f     func(done func()) // connect function
	delay int32             // delay in ms
}

// ============================================================================

func NewConnectQ() *ConnectQ {
	return &ConnectQ{
		q:    make(chan *connectF),
		quit: make(chan int),
	}
}

func (self *ConnectQ) Open() {
	self.wg.Add(1)
	go self.thrConnect()
}

func (self *ConnectQ) Close() {
	close(self.q)
	close(self.quit)
	self.wg.Wait()
}

func (self *ConnectQ) Connect(f func(done func()), delay int32) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	self.q <- &connectF{f, delay}
}

// ============================================================================

func (self *ConnectQ) thrConnect() {
	defer self.wg.Done()

	for cf := range self.q {
		cf := cf

		go func() {
			select {
			case <-self.quit:
				return

			case <-time.After(time.Duration(cf.delay) * time.Millisecond):
				self.wg.Add(1)
				cf.f(self.wg.Done)
			}
		}()
	}
}
