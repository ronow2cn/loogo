package asyncop

import (
	"comm/sched/loop"
	"sync"
)

// ============================================================================

var (
	q    = make(chan *asyncOPT, 100000)
	quit = make(chan int)
	wg   sync.WaitGroup
)

// ============================================================================

type asyncOPT struct {
	op func() // run in background thread
	cb func() // run in logic thread
}

// ============================================================================

func Start() {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-quit:
				return

			case aop := <-q:
				aop.op()
				if aop.cb != nil {
					loop.Push(aop.cb)
				}
			}
		}
	}()
}

func Stop() {
	close(quit)
	wg.Wait()
}

func Close() {
	close(q)
	for aop := range q {
		aop.op()
		if aop.cb != nil {
			aop.cb()
		}
	}
}

func Push(op func(), cb func()) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	q <- &asyncOPT{op, cb}
}
