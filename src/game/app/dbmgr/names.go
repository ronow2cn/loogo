package dbmgr

import (
	"comm/db"
)

// ============================================================================

// _id:  objectid
// name: string (uk)

// ============================================================================

func CenterInsertName(name string) bool {
	var rec struct {
		Name string `bson:"name"`
	}

	rec.Name = name

	err := DBCenter.Insert(CTabNameNames, &rec)
	return err == nil
}

func CenterChangeName(oldname, newname string) bool {
	err := DBCenter.UpdateByCond(
		CTabNameNames,
		db.M{"name": oldname},
		db.M{"$set": db.M{"name": newname}},
	)
	return err == nil
}
