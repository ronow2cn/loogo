package app

import (
	"comm"
	"comm/db"
	"comm/sched/asyncop"
	"comm/sched/loop"
	"game/app/dbmgr"
	"game/app/sched/resetter"
	"game/msg"
	"math/rand"
	"time"
)

// ============================================================================

type Player struct {
	user *User
	sid  uint64

	// save
	saveTs  time.Time
	saveTid *comm.Timer
}

// ============================================================================
// 核心功能

func newPlayer(u *User) *Player {
	return &Player{
		user:   u,
		saveTs: time.Now(),
	}
}

func (self *Player) open() {
	//user := self.user

}

func (self *Player) close() {
	// remove from resetter
	resetter.Remove(self)

	// stop save timer
	self.saveTimerStop()

	// final save
	self.Save()
}

func (self *Player) loaded() {

	// add to resetter
	resetter.Add(self)

	// start save timer
	self.saveTimerStart()
}

func (self *Player) created() {
	//user := self.user

	// 设置 初始配置

}

func (self *Player) OnOnline(ip string) {
	log.Info(self.user.Id, self.user.Name, "is online")

	// update login info
	now := time.Now()

	self.user.LoginTs = now
	self.user.LoginIP = ip

	asyncop.Push(
		func() {
			self.DB().Update(
				dbmgr.CTabNameUser,
				self.user.Id,
				db.M{
					"$set": db.M{
						"login_ts": now,
						"login_ip": ip,
					},
				},
			)
		},
		nil,
	)

	// update save timer: switch to online version
	self.saveTimerUpdate()
}

func (self *Player) OnOffline() {
	log.Info(self.user.Id, self.user.Name, "is offline")
}

func (self *Player) IsOnline() bool {
	return self.sid != 0
}

func (self *Player) DB() *db.Database {
	return self.user.db
}

func (self *Player) User() *User {
	return self.user
}

func (self *Player) SendMsg(message msg.Message) {
	if self.IsOnline() {
		NetMgr.Send2Player(self.sid, message)
	}
}

func (self *Player) Logout() {
	NetMgr.Send2Gate(sid2gateid(self.sid), &msg.GS_Kick{Sid: self.sid})
	PlayerMgr.SetOffline(self)
}

// ============================================================================
// 基础功能

func (self *Player) GetId() string {
	return self.user.Id
}

func (self *Player) GetName() string {
	return self.user.Name
}

func (self *Player) GetChannel() string {
	return self.user.Channel
}

func (self *Player) GetPlat() string {
	return self.user.Plat
}

func (self *Player) GetHead() int32 {
	return self.user.Head

}
func (self *Player) GetHFrame() int32 {
	return self.user.HFrame
}

func (self *Player) GetLevel() int32 {
	return self.user.Lv
}

func (self *Player) GetExp() int32 {
	return self.user.Exp
}

func (self *Player) ChangeName(name string, f func(bool)) {
	var b bool

	oldname := self.user.Name

	asyncop.Push(
		func() {
			// update name-db
			if !dbmgr.CenterChangeName(oldname, name) {
				b = false
				return
			}

			// update center-db
			dbmgr.CenterUpdateUserName(self.GetId(), name)

			// update game-db
			err := dbmgr.DBGame.Update(
				dbmgr.CTabNameUser,
				self.GetId(),
				db.M{"$set": db.M{"name": name}},
			)
			if err != nil {
				log.Warning("Player.ChangeName() failed:", err)
			}

			b = true
		},

		func() {
			// update memory
			if b {
				PlayerMgr.UpdatePlayerName(self, name)
			}

			// callback
			f(b)
		},
	)
}

// ============================================================================
// 存盘

func (self *Player) saveTimerStart() {
	self.saveTid = loop.SetTimeout(self.saveNextTime(), func() {
		self.SaveAsync()
		self.saveTimerStart()
	})
}

func (self *Player) saveTimerStop() {
	if self.saveTid != nil {
		loop.CancelTimer(self.saveTid)
		self.saveTid = nil
	}
}

func (self *Player) saveNextTime() time.Time {
	// 此方案是为了减轻玩家存盘压力

	if self.IsOnline() {
		// 在线存盘间隔
		return self.saveTs.Add(time.Duration(300+rand.Intn(300)) * time.Second)
	} else {
		// 离线存盘间隔
		return self.saveTs.Add(time.Duration(600+rand.Intn(600)) * time.Second)
	}
}

func (self *Player) saveTimerUpdate() {
	loop.UpdateTimer(self.saveTid, self.saveNextTime())
}

// 异步存盘：用于定时存盘
func (self *Player) SaveAsync() {
	// clone
	obj := CloneBsonObject(self.user)

	// async save
	asyncop.Push(
		func() {
			self.DB().Update(
				dbmgr.CTabNameUser,
				self.user.Id,
				obj,
			)
		},
		nil,
	)

	// update save time
	self.saveTs = time.Now()
}

// 同步存盘：用于停服时存盘
func (self *Player) Save() {
	self.DB().Update(
		dbmgr.CTabNameUser,
		self.user.Id,
		self.user,
	)

	// update save time
	self.saveTs = time.Now()
}

// ============================================================================
// 重置处理

func (self *Player) ResetGetTime() time.Time {
	return self.user.RstTs
}

func (self *Player) ResetSetTime(ts time.Time) {
	self.user.RstTs = ts
}

func (self *Player) ResetDaily() {

}
