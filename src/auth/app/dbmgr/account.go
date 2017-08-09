package dbmgr

import (
	"comm/db"
	"time"
)

// ============================================================================

type Account struct {
	Channel      int32     `bson:"channel"`     //账号渠道类型
	ChannelUid   string    `bson:"channel_uid"` //对应渠道UID
	Token        string    `bson:"token"`       //账号token
	ExpireT      time.Time `bson:"expire_t"`    //token过期时间
	RefreshToken string    `bson:"r_token"`     //刷新token(主要微信，支付宝用到)
}

// ============================================================================

func CenterGetAccountInfo(channel int32, uid string) *Account {
	var obj Account

	err := DBCenter.GetObjectByCond(
		CTabNameAccount,
		db.M{
			"channel":     channel,
			"channel_uid": uid,
		},
		&obj,
	)

	if err != nil {
		return nil
	}

	return &obj
}

func CenterUpdateAccountInfo(channel int32, uid string, token, refrtoken string, expire int64) {

	err := DBCenter.Upsert(
		CTabNameAccount,
		db.M{
			"channel":     channel,
			"channel_uid": uid,
		},
		&Account{
			Channel:      channel,
			ChannelUid:   uid,
			Token:        token,
			RefreshToken: refrtoken,
			ExpireT:      time.Now().Add(time.Duration(expire) * time.Second),
		},
	)

	if err != nil {
		log.Warning("CenterUpdateAccountInfo error", err)
		return
	}
}
