package dbmgr

import (
	"comm/config"
	"comm/db"
)

// ============================================================================

type userload struct {
	DbName string `bson:"_id"`
	N      int64  `bson:"n"`
}

// ============================================================================

func CenterCreateUserLoad() {
	var obj userload

	for k, _ := range config.Common.DBUser {
		obj.DbName = k
		obj.N = 0

		err := DBCenter.Insert(CTabNameUserload, &obj)
		if err != nil && !db.IsDup(err) {
			log.Error("dbmgr.CenterCreateUserLoad() failed:", err)
		}
	}
}
