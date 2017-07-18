package loop

import (
	"comm"
	"comm/logger"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================

var log = logger.DefaultLogger

// ============================================================================

var (
	q      = make(chan func(), 100000)
	timerq = comm.NewTimerQueue()
	quit   = make(chan int)
	wg     sync.WaitGroup

	numHandled int32
)

// ============================================================================

func Run() {
	wg.Add(1)

	go loopFunc()
	go loopTimer()
}

func Stop() {
	close(quit)
	wg.Wait()
}

func Push(f func()) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	q <- f
}

func SetTimeout(ts time.Time, f func()) *comm.Timer {
	return timerq.SetTimeout(ts, f)
}

func CancelTimer(t *comm.Timer) {
	timerq.Cancel(t)
}

func UpdateTimer(t *comm.Timer, ts time.Time) {
	timerq.Update(t, ts)
}

func QLen() int32 {
	return int32(len(q))
}

func NumHandled() int32 {
	return atomic.SwapInt32(&numHandled, 0)
}

// ============================================================================

func loopFunc() {
	defer wg.Done()

	for f := range q {
		safeExecute(f)
		atomic.AddInt32(&numHandled, 1)
	}
}

func loopTimer() {
	defer close(q)

	for {
		select {
		case <-quit:
			return

		default:
			Push(func() {
				now := time.Now()
				for timerq.Expire(now) {
				}
			})

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func safeExecute(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("critical exception:", err)
			log.Error(comm.Callstack())
		}
	}()

	f()
}
