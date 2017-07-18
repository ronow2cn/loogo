package resetter

import (
	"comm"
	"comm/sched/loop"
	"time"
)

// ============================================================================

const (
	CRESETHOUR = 3
)

// ============================================================================

type Resettable struct {
	RstTs time.Time // last reset timestamp
}

func (self *Resettable) ResetGetTime() time.Time {
	return self.RstTs
}

func (self *Resettable) ResetSetTime(ts time.Time) {
	self.RstTs = ts
}

type IResettable interface {
	ResetGetTime() time.Time
	ResetSetTime(ts time.Time)
	ResetDaily()
}

// ============================================================================

var (
	objs = make(map[IResettable]bool) // obj map
	tid  *comm.Timer                  // timer id
)

// ============================================================================

func Start() {
	if tid != nil {
		return
	}

	// sched for next reset
	key := calcLastestKeytime()
	schedDaily(key)
}

func Stop() {
	if tid == nil {
		return
	}

	loop.CancelTimer(tid)
	tid = nil
}

func Add(obj IResettable) {
	objs[obj] = true

	// check if we should reset immediately
	key := calcLastestKeytime()
	if obj.ResetGetTime().Before(key) {
		obj.ResetDaily()
		obj.ResetSetTime(key)
	}
}

func Remove(obj IResettable) {
	delete(objs, obj)
}

// ============================================================================

// return lastest key time before now
func calcLastestKeytime() time.Time {
	now := time.Now()

	y, M, d := now.Date()
	key := time.Date(y, M, d, CRESETHOUR, 0, 0, 0, time.Local)

	if key.After(now) {
		key = key.Add(-24 * time.Hour)
	}

	return key
}

func schedDaily(key time.Time) {
	key = key.Add(24 * time.Hour)
	tid = loop.SetTimeout(key, func() {
		resetDaily(key)
		schedDaily(key)
	})
}

func resetDaily(key time.Time) {
	for obj, _ := range objs {
		if obj.ResetGetTime().Before(key) {
			obj.ResetDaily()
			obj.ResetSetTime(key)
		}
	}
}
