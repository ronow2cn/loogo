package app

import (
	"comm"
	"comm/db"
	"game/app/dbmgr"
	"time"
)

// ============================================================================

type User struct {
	Id         string    `bson:"_id"`         // Id
	IdDefault  string    `bson:"id_default"`  // 默认 认证 Id
	IdFacebook string    `bson:"id_facebook"` // facebook Id
	Svr0       string    `bson:"svr0"`        // 原始服名称
	Svr        string    `bson:"svr"`         // 当前服名称
	Channel    string    `bson:"channel"`     // 渠道
	Plat       string    `bson:"plat"`        // 登陆平台
	Name       string    `bson:"name"`        // 名字
	CreateTs   time.Time `bson:"create_ts"`   // 创建时间
	LoginTs    time.Time `bson:"login_ts"`    // 上次登陆时间
	LoginIP    string    `bson:"login_ip"`    // 上次登录 IP

	RstTs time.Time `bson:"rst_ts"` // 上次重置时间

	Head   int32 `bson:"head"`   // 头像
	HFrame int32 `bson:"hframe"` // 相框
	Lv     int32 `bson:"lv"`     // 等级
	Exp    int32 `bson:"exp"`    // 经验

	db *db.Database `bson:"-"`
}

// ============================================================================

func createUser(uid string, f func(*User)) *User {
	// get user db
	dbname := getUserDBName(uid)
	udb := dbmgr.UserDB(dbname)
	if udb == nil {
		log.Error("get user db failed:", dbname)
		log.Error(comm.Callstack())
		return nil
	}

	// new user
	user := &User{}

	// init data
	user.Id = uid
	user.Name = GenRandName()
	user.CreateTs = time.Now()
	user.RstTs = time.Unix(0, 0)

	user.Head = 0
	user.HFrame = 0
	user.Lv = 1
	user.Exp = 0

	// --------------------------------

	// callback
	if f != nil {
		f(user)
	}

	// save to db
	err := udb.Insert(dbmgr.CTabNameUser, user)
	if err != nil {
		log.Error("create user failed:", uid, err)
		return nil
	}

	// // create mailbox in DB
	// mail.CreateMailBoxInDB(udb, user.Id)

	// // create chatbox in DB
	// chat.CreateChatBoxInDB(udb, user.Id)

	// update user name into center-db
	dbmgr.CenterUpdateUserName(user.Id, user.Name)

	// bind db
	user.db = udb

	// return
	return user
}
