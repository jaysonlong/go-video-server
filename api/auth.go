package main

import (
	"net/http"

	"github.com/midmis/go-video-server/api/defs"

	"github.com/midmis/go-video-server/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_USERNAME = "X-User-Name"

// 每次请求处理前校验session, 并填充user_name请求头,
// 后面的业务逻辑中只需要检查user_name请求头即可判断用户是否有效
func ValidateSession(r *http.Request) bool {
	sessionId := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sessionId) == 0 {
		// 这里使用Set不用Add，因为Set方法无论请求头是否存在都会写入，而Add只在不存在时写入
		r.Header.Set(HEADER_FIELD_USERNAME, "")
		return false
	}

	if loginName, expired := session.IsSessionExpired(sessionId); !expired {
		r.Header.Set(HEADER_FIELD_USERNAME, loginName)
		return true
	}

	r.Header.Set(HEADER_FIELD_USERNAME, "")
	return false
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	userName := r.Header.Get(HEADER_FIELD_USERNAME)
	if len(userName) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}

	return true
}
