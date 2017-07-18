package config

import (
	"comm"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
)

// ============================================================================

type commonT struct {
	Version  string            `json:"version"`
	VerMajor string            `json:"-"`
	VerMinor string            `json:"-"`
	VerBuild string            `json:"-"`
	LogLevel string            `json:"logLevel"`
	Perfmon  string            `json:"perfmon"`
	DBCenter string            `json:"dbCenter"`
	DBUser   map[string]string `json:"dbUser"`
}

type gateT struct {
	Id    int32  `json:"-"`
	Name  string `json:"-"`
	IPWan string `json:"ipWan"`
	Port  int32  `json:"port"`
}

type gameT struct {
	Id      int32  `json:"-"`
	Name    string `json:"-"`
	Addr4GW string `json:"addr4gw"`
	DBGame  string `json:"dbGame"`
	DBLog   string `json:"dbLog"`
}

type switcherT struct {
	IP    string `json:"ip"`
	Port  int32  `json:"port"`
	Token string `json:"token"`
}

type authT struct {
	IP    string `json:"ip"`
	Port  int32  `json:"port"`
	Token string `json:"token"`
}

// ============================================================================

type configT struct {
	Common   *commonT          `json:"common"`
	Switcher *switcherT        `json:"switcher"`
	Auth     *authT            `json:"auth"`
	Gates    map[string]*gateT `json:"gates"`
	Games    map[string]*gameT `json:"games"`
}

// ============================================================================

var (
	Common   *commonT
	Switcher *switcherT
	Auth     *authT

	Gates map[string]*gateT
	Games map[string]*gameT
)

var (
	DefaultGate *gateT
	DefaultGame *gameT
)

// ============================================================================

func Parse(fn string, server string) {
	var conf configT

	// read file
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		comm.Panic("open config file failed:", err)
	}

	// parse
	err = json.Unmarshal(d, &conf)
	if err != nil {
		comm.Panic("parse config file failed:", err)
	}

	parseIDName(&conf)
	parseVersion(&conf)

	// set variables
	Common = conf.Common
	Gates = conf.Gates
	Games = conf.Games
	Switcher = conf.Switcher
	Auth = conf.Auth

	// set defaults
	for k, v := range Gates {
		if k == server {
			DefaultGate = v
			DefaultGame = Games["game"+comm.I32toa(v.Id)]
			return
		}
	}

	for k, v := range Games {
		if k == server {
			DefaultGame = v
			DefaultGate = Gates["gate"+comm.I32toa(v.Id)]
			return
		}
	}
}

// ============================================================================

func parseIDName(conf *configT) {
	re := regexp.MustCompile(`^[a-z]+(\d+)$`)

	for k, v := range conf.Gates {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			comm.Panic("invalid gate name:", k)
		}

		v.Id = comm.Atoi32(arr[1])
		v.Name = k
	}

	for k, v := range conf.Games {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			comm.Panic("invalid game name:", k)
		}

		v.Id = comm.Atoi32(arr[1])
		v.Name = k
	}
}

func parseVersion(conf *configT) {
	arr := strings.Split(conf.Common.Version, ".")
	if len(arr) < 3 {
		comm.Panic("invalid version:", conf.Common.Version)
	}

	conf.Common.VerMajor = arr[0]
	conf.Common.VerMinor = arr[1]
	conf.Common.VerBuild = arr[2]
}
