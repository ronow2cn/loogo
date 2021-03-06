package dbmgr

import (
	"comm/config"
	"comm/db"
)

// ============================================================================

const (
	// center
	CTabNameUserinfo = "userinfo"
)

// ============================================================================

var (
	DBCenter *db.Database
)

// ============================================================================

func Open() {
	// init center db
	if DBCenter == nil {
		DBCenter = db.NewDatabase()
		DBCenter.Open(config.Common.DBCenter, true)
	}
}

func Close() {
	DBCenter.Close()
}
