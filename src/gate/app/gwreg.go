package app

import (
	"comm"
	"comm/config"
	"fmt"
	"time"
)

// ============================================================================

var (
	GWReg = &gwreg{}
)

// ============================================================================

type gwreg struct {
}

// ============================================================================

func (self *gwreg) Run() {
	go self.reg2switcher()
}

func (self *gwreg) reg2switcher() {
	for {
		time.Sleep(6 * time.Second)

		comm.HttpGet(
			fmt.Sprintf(
				"http://%s:%d/server/gatereg?token=%s&name=%s&ip=%s&port=%s",
				config.Switcher.IP,
				config.Switcher.Port,
				config.Switcher.Token,
				config.DefaultGame.Name,
				config.DefaultGate.IPWan,
				comm.I32toa(config.DefaultGate.Port),
			))
	}
}
