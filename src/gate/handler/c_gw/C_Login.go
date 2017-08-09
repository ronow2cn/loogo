package c_gw

import (
	"auth/route"
	"comm"
	"comm/config"
	"encoding/json"
	"fmt"
	"game/app/gconst"
	"gate/app"
	"gate/app/dbmgr"
	"gate/msg"
	"proto/errorcode"
	"proto/macrocode"
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
		if config.Games[req.Svr0] == nil {
			ec = Err.Login_Failed
			return
		}

		if req.AuthChannel == macrocode.ChannelType_Test {
			res.AuthId = req.AuthId
		} else {
			//auth channel uid
			authret := authAccount(req)
			if authret == nil || (authret.Result != Err.OK) {
				ec = Err.Login_Failed
				return
			}

			res.AuthId = authret.OpenId
			res.Token = authret.Token
			res.ExpireT = authret.Expire
			//================
		}
		// find user info
		user := dbmgr.CenterGetUserInfo(req.AuthChannel, res.AuthId, req.Svr0)
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

func authAccount(message msg.Message) *route.AuthRes {
	req := message.(*msg.C_Login)

	url := fmt.Sprintf("http://%s:%s/auth/%s?servertoken=%s&authtype=%s&uid=%s&token=%s",
		config.Auth.IP, comm.I32toa(config.Auth.Port), req.AuthChannel, config.Auth.Token, req.AuthType, req.AuthId, req.AuthToken)

	ret, err := comm.HttpGetT(url, gconst.HttpTimeOutSecond)
	if err != nil {
		log.Error("updateWeinXinTokenByRefreshToken HttpGetT error", err)
		return nil
	}

	var jret route.AuthRes
	err = json.Unmarshal([]byte(ret), &jret)
	if err != nil {
		log.Error("authAccount Unmarshal error", ret, err)
		return nil
	}

	return &jret
}
