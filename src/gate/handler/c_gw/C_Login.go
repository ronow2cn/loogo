package c_gw

import (
	//	"comm"
	//	"fmt"
	"comm/config"
	"gate/app"
	"gate/app/dbmgr"
	"gate/msg"
	"proto/errorcode"
	"time"
)

func C_Login(message msg.Message, ctx interface{}) {
	// !Note: in net-thread

	req := message.(*msg.C_Login)
	sess := ctx.(*app.Session)

	ec := Err.OK
	res := &msg.GW_Login_R{}

	go func() {
		defer func() {
			// end auth
			sess.EndAuth()
			res.ErrorCode = ec
			res.AuthId = req.AuthId

			sess.SendMsg(res)
		}()

		// check version
		if req.VerMajor != config.Common.VerMajor || req.VerMinor != config.Common.VerMinor {
			ec = Err.Login_InvalidVersion
			return
		}

		if !sess.BeginAuth() {
			ec = Err.Login_Failed
			return
		}

		// check param
		if req.AuthId == "" {
			ec = Err.Login_Failed
			return
		}

		if config.Games[req.Svr0] == nil {
			ec = Err.Login_Failed
			return
		}

		//auth channel uid

		//================

		// find user info
		user := dbmgr.CenterGetUserInfo(req.AuthChannel, req.AuthId, req.Svr0)
		if user == nil {
			ec = Err.Login_Failed
			return
		}

		// check is in ban time
		if user.BanTs.After(time.Now()) {
			ec = Err.Login_UserBanned
			return
		}

		// login player: notify gamesvr
		sess.LoginPlayer(&msg.GW_UserOnline{
			Sid:        sess.GetId(),
			UserId:     user.UserId,
			Channel:    user.Channel,
			ChannelUid: user.ChannelUid,
			Svr0:       user.Svr0,
			LoginIP:    sess.GetIP(),
		})

	}()

}
