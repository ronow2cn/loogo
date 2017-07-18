package dbmgr

import (
	"comm/db"
	"fmt"
)

// ============================================================================

type seqidT struct {
	Id         int   `bson:"_id"`
	SeqUserId  int64 `bson:"seq_uid"`
	SeqGuildId int64 `bson:"seq_gid"`
}

// ============================================================================

func CenterGenUserId(dbname string) string {
	var obj seqidT

	err := DBCenter.FindAndModify(
		CTabNameSeqid,
		db.M{"_id": 1},
		db.Change{
			Update: db.M{
				"$inc": db.M{"seq_uid": 1},
			},
			ReturnNew: true,
		},
		db.M{"seq_uid": 1},
		&obj,
	)
	if err != nil {
		log.Error("dbmgr.CenterGenUserId() failed:", err)
	}

	return fmt.Sprintf("%s-%d", dbname, obj.SeqUserId)
}
