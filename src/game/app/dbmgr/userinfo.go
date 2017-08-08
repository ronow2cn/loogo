package dbmgr

import (
	"comm/db"
	"time"
)

// ============================================================================

type userinfo struct {
	UserId     string    `bson:"_id"`         //游戏中玩家唯一id
	Channel    int32     `bson:"channel"`     //账号渠道类型
	ChannelUid string    `bson:"channel_uid"` //对应渠道UID
	Svr0       string    `bson:"svr0"`        //玩家原始服务器
	Svr        string    `bson:"svr"`         //玩家登陆服务器
	Name       string    `bson:"name"`        //玩家名字
	BanTs      time.Time `bson:"ban_ts"`      //玩家封号时间
}

// ============================================================================

func CenterUpdateUserName(userid string, name string) {
	err := DBCenter.Update(
		CTabNameUserinfo,
		userid,
		db.M{"$set": db.M{"name": name}},
	)
	if err != nil {
		log.Error("dbmgr.CenterUpdateUserName() failed:", err)
	}
}

func CenterBanUser(userid string, min int32) {
	// > 0: 封号
	// < 0: 解封

	if min == 0 {
		return
	}

	bants := time.Unix(0, 0)
	if min > 0 {
		bants = time.Now().Add(time.Duration(min) * time.Minute)
	}

	err := DBCenter.Update(
		CTabNameUserinfo,
		userid,
		db.M{"$set": db.M{"bants": bants}},
	)
	if err != nil {
		log.Error("dbmgr.CenterBanUser() failed:", err)
	}
}

func CenterGetBanInfo(userid string) time.Time {
	var obj userinfo

	err := DBCenter.GetProjectionByCond(
		CTabNameUserinfo,
		db.M{"_id": userid},
		db.M{"ban_ts": 1},
		&obj,
	)
	if err != nil {
		log.Error("dbmgr.CenterGetBanInfo() failed:", err)
		return time.Unix(0, 0)
	}

	return obj.BanTs
}
