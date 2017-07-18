package dbmgr

import (
	"comm/db"
	"gopkg.in/mgo.v2"
)

// ============================================================================

type userloadT struct {
	DbName string `bson:"_id"`
	N      int64  `bson:"n"`
}

// ============================================================================

func CenterAllocUserDB() string {
	var obj userloadT
	var err error

	DBCenter.Execute(func(session *mgo.Session) {
		err = session.DB("").C(CTabNameUserload).Find(nil).Sort("n").Limit(1).One(&obj)
	})

	if err == nil {
		// ok
		return obj.DbName
	} else {
		// failed
		log.Error("dbmgr.CenterAllocUserDB() failed:", err)
		return ""
	}
}

func CenterIncUserLoad(dbname string) {
	err := DBCenter.Update(
		CTabNameUserload,
		dbname,
		db.M{"$inc": db.M{"n": 1}},
	)
	if err != nil {
		log.Error("dbmgr.CenterIncUserLoad() failed:", err)
	}
}

func CenterDecUserLoad(dbname string) {
	err := DBCenter.Update(
		CTabNameUserload,
		dbname,
		db.M{"$inc": db.M{"n": -1}},
	)
	if err != nil {
		log.Error("dbmgr.CenterDecUserLoad() failed:", err)
	}
}
