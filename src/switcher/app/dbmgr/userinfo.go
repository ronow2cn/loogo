package dbmgr

import (
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

func CenterGetUserSvr(channel int32, authid string, svr0 string) string {
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
		return obj.Svr

	} else {
		return svr0
	}
}
