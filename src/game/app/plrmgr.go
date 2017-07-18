package app

import (
	"comm"
	"comm/config"
	"comm/db"
	"game/app/dbmgr"
	"strings"
	"sync/atomic"
)

// ============================================================================

var PlayerMgr = &plrmgrT{
	plrsByID:   make(map[string]*Player),
	plrsByName: make(map[string]*Player),
	plrsOnline: make(map[uint64]*Player),
}

// ============================================================================

type plrmgrT struct {
	plrsByID   map[string]*Player // loaded players by id
	plrsByName map[string]*Player // loaded players by name
	plrsOnline map[uint64]*Player // online players by sid

	numLoaded int32 // loaded player number
	numOnline int32 // online player number
}

// ============================================================================

func getUserDBName(uid string) string {
	return strings.Split(uid, "-")[0]
}

// ============================================================================

func (self *plrmgrT) Close() {
	for _, plr := range self.plrsByID {
		plr.close()
	}

	log.Info("ALL players are saved")
}

func (self *plrmgrT) LoadPlayer(uid string, ignoreErr ...bool) *Player {
	//check
	if uid == "" {
		return nil
	}

	// find in memory
	plr := self.plrsByID[uid]
	if plr == nil {
		// load user from db
		user := self.loadFromDB(uid)
		if user == nil {
			if len(ignoreErr) == 0 || !ignoreErr[0] {
				log.Error("user not found:", uid)
				log.Error(comm.Callstack())
			}
			return nil
		}

		// create player
		plr = newPlayer(user)

		// add to mgr
		self.addPlayer(plr, false)
	}

	return plr
}

func (self *plrmgrT) CreatePlayer(uid string, f func(*User)) *Player {
	user := createUser(uid, f)
	if user == nil {
		return nil
	}

	// create player
	plr := newPlayer(user)

	// add to mgr
	self.addPlayer(plr, true)

	return plr
}

func (self *plrmgrT) SetOnline(plr *Player, sid uint64, ip string) {
	if plr.sid != 0 || sid == 0 {
		return
	}

	// set sid
	plr.sid = sid

	// add to mgr
	self.plrsOnline[sid] = plr

	// event: online
	plr.OnOnline(ip)

	// count
	atomic.AddInt32(&self.numOnline, 1)
}

func (self *plrmgrT) SetOffline(plr *Player) {
	if plr.sid == 0 {
		return
	}

	// remove from mgr
	delete(self.plrsOnline, plr.sid)

	// reset sid
	plr.sid = 0

	// event: offline
	plr.OnOffline()

	// count
	atomic.AddInt32(&self.numOnline, -1)
}

func (self *plrmgrT) OfflineAllPlayers(gateid int32) {
	if gateid == 0 {
		// all players
		for _, plr := range self.plrsByID {
			self.SetOffline(plr)
		}
	} else {
		// players from gateid
		for _, plr := range self.plrsByID {
			if sid2gateid(plr.sid) == gateid {
				self.SetOffline(plr)
			}
		}
	}
}

func (self *plrmgrT) FindPlayerById(uid string) *Player {
	return self.plrsByID[uid]
}

func (self *plrmgrT) FindPlayerByName(name string) *Player {
	return self.plrsByName[name]
}

func (self *plrmgrT) FindPlayerBySid(sid uint64) *Player {
	return self.plrsOnline[sid]
}

func (self *plrmgrT) ArrayLoadedPlayers() (ret []*Player) {
	ret = make([]*Player, 0, len(self.plrsByID))
	for _, plr := range self.plrsByID {
		ret = append(ret, plr)
	}
	return
}

func (self *plrmgrT) ArrayOnlinePlayers() (ret []*Player) {
	ret = make([]*Player, 0, len(self.plrsOnline))
	for _, plr := range self.plrsOnline {
		ret = append(ret, plr)
	}
	return
}

func (self *plrmgrT) UpdatePlayerName(plr *Player, name string) {
	if name == plr.user.Name {
		return
	}

	delete(self.plrsByName, plr.user.Name)
	plr.user.Name = name
	self.plrsByName[name] = plr
}

func (self *plrmgrT) NumLoaded() int32 {
	return atomic.LoadInt32(&self.numLoaded)
}

func (self *plrmgrT) NumOnline() int32 {
	return atomic.LoadInt32(&self.numOnline)
}

// ============================================================================

func (self *plrmgrT) loadFromDB(uid string) *User {
	// get user db
	dbname := getUserDBName(uid)
	udb := dbmgr.UserDB(dbname)
	if udb == nil {
		log.Critical("get user db failed:", dbname)
		log.Critical(comm.Callstack())
		return nil
	}

	// load
	var user User
	err := udb.GetObjectByCond(
		dbmgr.CTabNameUser,
		db.M{
			"_id": uid,
			"svr": config.DefaultGame.Name,
		},
		&user,
	)
	if err != nil {
		return nil
	}

	// bind db
	user.db = udb

	// return
	return &user
}

func (self *plrmgrT) addPlayer(plr *Player, creation bool) {
	self.plrsByID[plr.user.Id] = plr
	self.plrsByName[plr.user.Name] = plr

	// open
	plr.open()

	// loaded
	plr.loaded()

	// created
	if creation {
		plr.created()
	}

	// count
	atomic.AddInt32(&self.numLoaded, 1)
}
