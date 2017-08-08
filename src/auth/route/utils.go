package route

import (
	"comm/logger"
)

var log = logger.DefaultLogger

func checkServerToken(token string) bool {
	return len(token) > 0 && token == config.Auth.Token
}