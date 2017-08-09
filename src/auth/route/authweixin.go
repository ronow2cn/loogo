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
	authType := r.PostFormValue("authtype")
	openid := r.PostFormValue("uid")

	if authType != comm.I32toa(macrocode.LoginType_WeiXinCode) { //token auth
		processWeiXinTokenAuth(w, openid, token)
	} else { //code auth
		processWeiXinCodeAuth(w, token)
	}
}

func processWeiXinTokenAuth(w http.ResponseWriter, openid string, token string) {

}

func processWeiXinCodeAuth(w http.ResponseWriter, code string) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		config.Auth.WeiXin.TokenUrl, config.Auth.WeiXin.AppId, config.Auth.WeiXin.AppKey, code)

	ret, err := comm.HttpGetT(url, HttpTimeOutSecond)
	if err != nil {
		log.Error("processWeiXinCodeAuth HttpGetT error", err)
		makeAuthRes(w, &AuthRes{Result: Err.Failed})
		return
	}

	var jret weixinTokenRet
	err = json.Unmarshal([]byte(ret), &jret)
	if err != nil {
		log.Error("processWeiXinCodeAuth Unmarshal error", err)
		makeAuthRes(w, &AuthRes{Result: Err.Failed})
		return
	}

	if len(jret.ErrCode) != 0 {
		log.Error("processWeiXinCodeAuth ErrCode", jret.ErrCode, jret.ErrMsg)
		makeAuthRes(w, &AuthRes{Result: comm.Atoi32(jret.ErrCode)})
		return
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
}

func makeAuthRes(w http.ResponseWriter, ret *AuthRes) {
	jres, err := json.Marshal(ret)
	if err != nil {
		log.Error("makeAuthRes Marshal failed")
		return
	}

	fmt.Fprint(w, string(jres))
}
