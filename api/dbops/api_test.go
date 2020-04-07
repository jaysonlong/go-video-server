package dbops

import (
	"testing"
	"time"

	"github.com/midmis/go-video-server/api/defs"
)

var (
	testUserName       = "zhangsan"
	testPwd            = "123"
	testVideoName      = "冰雪奇缘"
	tmpVideoId         string
	testVideoId        = "abcd-1234"
	testCommentContent = "this is a good video"
	testSessionId      = []string{"abcd", "efgh"}
	testUserInfo       *defs.UserInfo
)

func clearTables() {
	dbConn.Exec("TRUNCATE TABLE users")
	dbConn.Exec("TRUNCATE TABLE video_info")
	dbConn.Exec("TRUNCATE TABLE comments")
	dbConn.Exec("TRUNCATE TABLE sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

// test user api
func TestUserWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUserInfo)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUser(testUserName, testPwd)
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUserInfo(t *testing.T) {
	userInfo, err := GetUserInfo(testUserName)
	if err != nil {
		t.Errorf("Error of GetUserInfo: %v", err)
	}

	if userInfo.Pwd != testPwd {
		t.Errorf("pwd is %v, not equal %v", userInfo.Pwd, testPwd)
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser(testUserName, testPwd)
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}
}

func testRegetUser(t *testing.T) {
	userInfo, err := GetUserInfo(testUserName)
	if err != nil {
		t.Errorf("Error of RegetUser: %v", err)
	}

	if userInfo != nil {
		t.Errorf("userInfo is %v, not nil", userInfo)
	}
}

// test video_info api
func TestVideoInfoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Add", testAddVideoInfo)
	t.Run("Get", testGetVideoInfo)
	t.Run("List", testListVideoInfo)
	t.Run("Del", testDeleteVideoInfo)
	t.Run("Reget", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	AddUser(testUserName, testPwd)
	userInfo, err := GetUserInfo(testUserName)
	testUserInfo = userInfo
	if err != nil {
		t.Errorf("Prepare user info failed: %v", err)
	}

	videoInfo, err := AddVideoInfo(testVideoName, testUserInfo.UserId)
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}

	if videoInfo == nil {
		t.Errorf("videoInfo is nil")
	}
	tmpVideoId = videoInfo.VideoId
}

func testGetVideoInfo(t *testing.T) {
	videoInfo, err := GetVideoInfo(tmpVideoId)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}

	if videoInfo == nil || videoInfo.Name != testVideoName {
		t.Errorf("videoInfo.Name is %v, not equal %v", videoInfo.Name, testVideoName)
	}

}

func testListVideoInfo(t *testing.T) {
	AddVideoInfo(testVideoName+"_ex", testUserInfo.UserId)
	currentTS := time.Now().Unix()
	videoInfoArr, err := ListVideoInfo(testUserName, currentTS-1000*3600*24, currentTS)
	if err != nil {
		t.Errorf("Error of ListVideoInfo: %v", err)
	}

	if len(videoInfoArr) != 2 {
		t.Errorf("len(videoInfoArr) is %v, not equal %v", len(videoInfoArr), 2)
	}

}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tmpVideoId)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	videoInfo, err := GetVideoInfo(tmpVideoId)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}

	if videoInfo != nil {
		t.Errorf("videoInfo is not nil: %v", videoInfo)
	}
}

// test comment api
func TestCommentWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Add", testAddComment)
	t.Run("Add", testAddComment)
	t.Run("List", testListComments)
}

func testAddComment(t *testing.T) {
	comment, err := AddComment(testCommentContent, 1, testVideoId)
	if err != nil {
		t.Errorf("Error of AddComment: %v", err)
	}

	if comment == nil {
		t.Errorf("comment is nil")
	}
}

func testListComments(t *testing.T) {
	currentTS := time.Now().Unix()
	comments, err := ListComments(testVideoId, currentTS-1000*3600*24, currentTS)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("comment count not equal 2")
	}
	if comments[0].Content != testCommentContent {
		t.Errorf("the last comment is %v, not equal %v", comments[0].Content, testCommentContent)
	}

}

// test session api
func TestSessionWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Add", testAddSession)
	t.Run("Retrieve", testRetrieveSession)
	t.Run("RetrieveAll", testRetrieveAllSession)
	t.Run("Delete", testDeleteSession)
	t.Run("Re-Retrieve", testReRetrieveSession)
}

func testAddSession(t *testing.T) {
	currentTime := int(time.Now().Unix())
	err := AddSession(testSessionId[0], testUserName, currentTime+60*30)
	if err != nil {
		t.Errorf("Error of AddSession: %v", err)
	}

	err = AddSession(testSessionId[1], testUserName, currentTime+60*30)
	if err != nil {
		t.Errorf("Error of AddSession: %v", err)
	}
}

func testRetrieveSession(t *testing.T) {
	session, err := RetrieveSession(testSessionId[0])
	if err != nil {
		t.Errorf("Error of RetrieveSession: %v", err)
	}

	if session.LoginName != testUserName {
		t.Errorf("the session login name is %v, not equal %v", session.LoginName, testUserName)
	}

}

func testRetrieveAllSession(t *testing.T) {
	sessionMap, err := RetrieveAllSession()
	if err != nil {
		t.Errorf("Error of RetrieveAllSession: %v", err)
	}

	var cnt int
	sessionMap.Range(func(key, value interface{}) bool {
		cnt++
		return true
	})

	if cnt != 2 {
		t.Errorf("session count not equal 2")
	}
}

func testDeleteSession(t *testing.T) {
	err := DeleteSession(testSessionId[0])
	if err != nil {
		t.Errorf("Error of DeleteSession: %v", err)
	}
}

func testReRetrieveSession(t *testing.T) {
	session, err := RetrieveSession(testSessionId[0])
	if err != nil {
		t.Errorf("Error of RetrieveSession: %v", err)
	}

	if session != nil {
		t.Errorf("the session is not nil: %v", session)
	}

}
