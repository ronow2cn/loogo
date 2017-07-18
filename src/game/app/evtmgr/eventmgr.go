package evtmgr

import (
	"comm/logger"
)

// ============================================================================

var log = logger.DefaultLogger

// ============================================================================

var (
	evtMap = make(map[string][]func(...interface{}))
)

// ============================================================================
//register event, usage: func On() put in package func init()
func On(evt string, f func(...interface{})) {
	evtMap[evt] = append(evtMap[evt], f)
}

//call event
func Fire(evt string, args ...interface{}) {
	for _, f := range evtMap[evt] {
		f(args...)
	}
}
