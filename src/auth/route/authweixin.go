package route

import (
	"auth/app/dbmgr"
	"comm"
	"comm/config"
	"comm/sched/asyncop"
	"encoding/json"
	"fmt"
	"net/http"
	"proto/errorcode"
	"proto/macrocode"
	"time"
)

// ============================================================================
//weixin server res
type weixinTokenRet struct {
	Token     string `json:"access_token"`
	ExpireIn  int64  `json:"expires_in"`
	RefrToken string `json:"refresh_token"`
	OpenId    string `json:"openid"`
	Scope     string `json:"scope"`
	UnionId   string `json:"unionid"`
	ErrCode   string `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

// ============================================================================

func HandlerWeiXinAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method != "POST" {
		return
	}

	//werver token
	serverToken := r.PostFormValue("servertoken")
	if !checkServerToken(serverToken) {
		makeAuthRes(w, &AuthRes{Result: Err.Failed})
		return
	}

	//auth token
	token := r.PostFormValue("token")
	if len(token) == 0 {
		makeAuthRes(w, &AuthRes{Result: Err.Failed})
		return
	}

	//auth type: default token auth
	err := Err.OK
	authType := r.PostFormValue("authtype")
	if authType != comm.I32toa(macrocode.LoginType_WeiXinCode) { //token auth
		openid := r.PostFormValue("uid")
		err = processWeiXinTokenAuth(w, openid, token)
	} else { //code auth
		err = processWeiXinCodeAuth(w, token)
	}

	if err != Err.OK {
		makeAuthRes(w, &AuthRes{Result: int32(err)})
	}
}

func processWeiXinTokenAuth(w http.ResponseWriter, openid string, token string) int {
	account := dbmgr.CenterGetAccountInfo(macrocode.ChannelType_WeiXin, openid)
	if account == nil {
		log.Error("processWeiXinTokenAuth get account failed")
		return Err.Failed
	}

	if token != account.Token {
		log.Error("processWeiXinTokenAuth token error")
		return Err.Failed
	}

	if time.Now().After(account.ExpireT) {
		err := updateWeinXinTokenByRefreshToken(w, account.RefreshToken)
		return err
	} else {
		makeAuthRes(w, &AuthRes{
			Result: Err.OK,
			OpenId: openid,
			Token:  token,
			Expire: account.ExpireT.Unix(),
		})
	}

	return Err.OK
}

func updateWeinXinTokenByRefreshToken(w http.ResponseWriter, refrtoken string) int {
	url := fmt.Sprintf("%s?appid=%s&grant_type=refresh_token&refresh_token=%s",
		config.Auth.WeiXin.RefrUrl, config.Auth.WeiXin.AppId, refrtoken)

	ret, err := comm.HttpGetT(url, HttpTimeOutSecond)
	if err != nil {
		log.Error("updateWeinXinTokenByRefreshToken HttpGetT error", err)
		return Err.Failed
	}

	var jret weixinTokenRet
	err = json.Unmarshal([]byte(ret), &jret)
	if err != nil {
		log.Error("processWeiXinCodeAuth Unmarshal error", err)
		return Err.Failed
	}

	if len(jret.ErrCode) != 0 {
		log.Error("processWeiXinCodeAuth ErrCode", jret.ErrCode, jret.ErrMsg)
		return Err.Failed
	}

	makeAuthRes(w, &AuthRes{
		Result: Err.OK,
		OpenId: jret.OpenId,
		Token:  jret.Token,
		Expire: jret.ExpireIn + time.Now().Unix(),
	})

	//flush
	asyncop.Push(func() {
		dbmgr.CenterUpdateAccountInfo(macrocode.ChannelType_WeiXin, jret.OpenId, jret.Token, jret.RefrToken, jret.ExpireIn)
	}, nil)

	return Err.OK
}

func processWeiXinCodeAuth(w http.ResponseWriter, code string) int {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		config.Auth.WeiXin.TokenUrl, config.Auth.WeiXin.AppId, config.Auth.WeiXin.AppKey, code)

	ret, err := comm.HttpGetT(url, HttpTimeOutSecond)
	if err != nil {
		log.Error("processWeiXinCodeAuth HttpGetT error", err)
		return Err.Failed
	}

	var jret weixinTokenRet
	err = json.Unmarshal([]byte(ret), &jret)
	if err != nil {
		log.Error("processWeiXinCodeAuth Unmarshal error", err)
		return Err.Failed
	}

	if len(jret.ErrCode) != 0 {
		log.Error("processWeiXinCodeAuth ErrCode", jret.ErrCode, jret.ErrMsg)
		return Err.Failed
	}

	makeAuthRes(w, &AuthRes{
		Result: Err.OK,
		OpenId: jret.OpenId,
		Token:  jret.Token,
		Expire: jret.ExpireIn + time.Now().Unix(),
	})

	//flush
	asyncop.Push(func() {
		dbmgr.CenterUpdateAccountInfo(macrocode.ChannelType_WeiXin, jret.OpenId, jret.Token, jret.RefrToken, jret.ExpireIn)
	}, nil)

	return Err.OK
}

func makeAuthRes(w http.ResponseWriter, ret *AuthRes) {
	jres, err := json.Marshal(ret)
	if err != nil {
		log.Error("makeAuthRes Marshal failed")
		return
	}

	fmt.Fprint(w, string(jres))
}
