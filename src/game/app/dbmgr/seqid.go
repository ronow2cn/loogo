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

func CenterCreateSeqId() {
	if DBCenter.HasCollection(CTabNameSeqid) {
		return
	}

	var obj seqidT

	obj.Id = 1
	obj.SeqUserId = 999999
	obj.SeqGuildId = 999999

	err := DBCenter.Insert(CTabNameSeqid, &obj)
	if err != nil {
		log.Error("dbmgr.CenterCreateSeqId() failed:", err)
	}
}

func CenterGenGuildId() string {
	var obj seqidT

	err := DBCenter.FindAndModify(
		CTabNameSeqid,
		db.M{"_id": 1},
		db.Change{
			Update: db.M{
				"$inc": db.M{"seq_gid": 1},
			},
			ReturnNew: true,
		},
		db.M{"seq_gid": 1},
		&obj,
	)
	if err != nil {
		log.Error("dbmgr.CenterGenGuildId() failed:", err)
	}

	return fmt.Sprintf("g-%d", obj.SeqGuildId)
}
