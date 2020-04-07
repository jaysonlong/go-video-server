package dbops

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/midmis/go-video-server/api/defs"
)

var httpClient = &http.Client{}

func AddSession(sessionId string, userName string, ttl int) error {
	stmt, err := dbConn.Prepare("INSERT INTO sessions (session_id, login_name, TTL) VALUES(?, ?, ?)")
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	_, err = stmt.Exec(sessionId, userName, ttl)
	defer stmt.Close()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func RetrieveSession(sessionId string) (*defs.SimpleSession, error) {
	stmt, err := dbConn.Prepare("SELECT login_name, TTL FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var (
		loginName string
		ttl       int
	)
	err = stmt.QueryRow(sessionId).Scan(&loginName, &ttl)
	defer stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%v", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	session := &defs.SimpleSession{SessionId: sessionId, LoginName: loginName, TTL: ttl}

	return session, nil
}

func RetrieveAllSession() (*sync.Map, error) {
	stmt, err := dbConn.Prepare("SELECT session_id, login_name, TTL FROM sessions")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	rows, err := stmt.Query()
	defer stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%v", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var (
		sessionMap = &sync.Map{}
		sessionId  string
		loginName  string
		ttl        int
	)
	for rows.Next() {
		if err = rows.Scan(&sessionId, &loginName, &ttl); err != nil {
			log.Printf("%v", err)
			break
		}

		session := &defs.SimpleSession{SessionId: sessionId, LoginName: loginName, TTL: ttl}
		sessionMap.Store(sessionId, session)
	}

	return sessionMap, nil
}

func DeleteSession(sessionId string) error {
	stmt, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	_, err = stmt.Exec(sessionId)
	defer stmt.Close()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func RequestDeleteVideoFile(videoId string) {
	url := defs.VIDEO_DEL_SERVER_HOST + "/video-delete-record/" + videoId
	req, _ := http.NewRequest("GET", url, nil)
	httpClient.Do(req)
}
