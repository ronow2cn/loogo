package app

import (
	"comm"
	"comm/config"
	"comm/sched/loop"
	"time"
)

var Gates = map[string]*gate{}

type gate struct {
	IP   string
	Port string
	Ts   int64 //last reg time
}

// ============================================================================

func RegGate(name, ip, port string) {
	_, ok := config.Games[name]
	if !ok {
		log.Error("RegGate invalid param", name)
		return
	}

	Gates[name] = &gate{
		IP:   ip,
		Port: port,
		Ts:   time.Now().Unix(),
	}

	log.Info("RegGate:", name, ip, port)
	return
}

func UpdateGates() {
	for k, v := range Gates {
		if time.Now().Unix() > v.Ts+60 {
			delete(Gates, k)
		}
	}
}

func PrintGates() {
	for k, v := range Gates {
		log.Info(k, v.IP, v.Port, v.Ts)
	}
}

// ============================================================================
//tick gates
var GatesTimer = &gateTimer{}

type gateTimer struct {
	tid *comm.Timer
}

func (self *gateTimer) Init() {
	self.update()
}

func (self *gateTimer) update() {
	self.tid = loop.SetTimeout(self.tickNextTime(), func() {
		UpdateGates()
		self.update()
	})
}
func (self *gateTimer) tickNextTime() time.Time {
	return time.Now().Add(time.Duration(60) * time.Second)
}
