package session

import (
	"sync"
	"time"

	"github.com/midmis/go-video-server/api/dbops"
	"github.com/midmis/go-video-server/api/defs"
	"github.com/midmis/go-video-server/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionsFromDB() {
	sessions, err := dbops.RetrieveAllSession()
	if err != nil {
		panic(err.Error())
	}

	sessions.Range(func(sessionId, session interface{}) bool {
		session = session.(*defs.SimpleSession)
		sessionMap.Store(sessionId, session)
		return true
	})
}

func GenerateNewSession(userName string) string {
	sessionId, _ := utils.NewUUID()
	ttl := int(time.Now().Unix() + 60*30)
	err := dbops.AddSession(sessionId, userName, int(ttl))
	if err != nil {
		return ""
	}

	session := &defs.SimpleSession{SessionId: sessionId, LoginName: userName, TTL: ttl}
	sessionMap.Store(sessionId, session)
	return sessionId
}

func DeleteSession(sessionId string) {
	sessionMap.Delete(sessionId)
	dbops.DeleteSession(sessionId)
}

func IsSessionExpired(sessionId string) (string, bool) {
	ss, ok := sessionMap.Load(sessionId)
	if ok {
		currentTime := int(time.Now().Unix())
		session := ss.(*defs.SimpleSession)

		if session.TTL < currentTime {
			DeleteSession(sessionId)
			return "", true
		}

		return session.LoginName, false
	}

	return "", true
}
