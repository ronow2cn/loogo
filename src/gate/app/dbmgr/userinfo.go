package dbmgr

import (
	"comm/config"
	"comm/db"
	"time"
)

// ============================================================================

type UserInfo struct {
	UserId     string    `bson:"_id"`         //游戏中玩家唯一id
	Channel    int32     `bson:"channel"`     //账号渠道类型
	ChannelUid string    `bson:"channel_uid"` //对应渠道UID
	Svr0       string    `bson:"svr0"`        //玩家原始服务器
	Svr        string    `bson:"svr"`         //玩家登陆服务器
	Name       string    `bson:"name"`        //玩家名字
	BanTs      time.Time `bson:"ban_ts"`      //玩家封号时间
}

// ============================================================================

func CenterGetUserInfo(channel int32, authid string, svr0 string) *UserInfo {
	var obj UserInfo

	err := DBCenter.GetObjectByCond(
		CTabNameUserinfo,
		db.M{
			"channel":     channel,
			"channel_uid": authid,
			"svr0":        svr0,
		},
		&obj,
	)
	if err == nil {
		// userinfo exists
		// check svr
		if obj.Svr == config.DefaultGame.Name {
			return &obj
		} else {
			return nil
		}
	} else if db.IsNotFound(err) {
		// allocate user db
		dbname := CenterAllocUserDB()
		if dbname == "" {
			return nil
		}

		obj.UserId = CenterGenUserId(dbname)
		obj.ChannelUid = authid
		obj.Svr0 = svr0
		obj.Svr = config.DefaultGame.Name
		obj.BanTs = time.Unix(0, 0)
		obj.Channel = channel

		// flush to db
		err = DBCenter.Insert(CTabNameUserinfo, &obj)
		if err == nil {
			// update user load
			CenterIncUserLoad(dbname)

			// return new userinfo
			return &obj
		} else {
			// failed
			return nil
		}
	} else {
		// failed
		return nil
	}
}
