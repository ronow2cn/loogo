package route

import (
	"comm/config"
	"comm/logger"
)

var log = logger.DefaultLogger

const (
	HttpTimeOutSecond = 5
)

// ============================================================================
//res to game server by http
type AuthRes struct {
	Result int32  `json:"result"`
	OpenId string `json:"openid"`
	Token  string `json:"token"`
	Expire int64  `json:"expire"` //token expire time
}

// ============================================================================

func checkServerToken(token string) bool {
	return len(token) > 0 && token == config.Auth.Token
}
