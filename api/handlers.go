package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/midmis/go-video-server/api/dbops"
	"github.com/midmis/go-video-server/api/defs"
	"github.com/midmis/go-video-server/api/session"
)

func createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, _ := ioutil.ReadAll(r.Body)
	credential := defs.UserCredential{}

	if err := json.Unmarshal(body, &credential); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUser(credential.UserName, credential.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	sessionId := session.GenerateNewSession(credential.UserName)
	resp := defs.SignUpResp{Success: true, SessionId: sessionId}

	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 201, string(res))
	}
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userName := ps.ByName("user_name")
	body, _ := ioutil.ReadAll(r.Body)
	credential := defs.UserCredential{UserName: userName}

	if err := json.Unmarshal(body, &credential); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	userInfo, err := dbops.GetUserInfo(credential.UserName)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	if userInfo == nil || userInfo.Pwd == "" || userInfo.Pwd != credential.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	sessionId := session.GenerateNewSession(credential.UserName)
	resp := defs.LoginResp{Success: true, SessionId: sessionId}

	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 200, string(res))
	}
}

func getUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	userName := ps.ByName("user_name")
	userInfo, err := dbops.GetUserInfo(userName)
	if err != nil {
		log.Printf("Erorr in GetUserinfo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	resp := &defs.UserInfoResp{UserId: userInfo.UserId, UserName: userName}
	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 200, string(res))
	}
}

func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	sessionId := r.Header.Get(HEADER_FIELD_SESSION)
	session.DeleteSession(sessionId)
	resp := defs.LogoutResp{Success: true}

	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 200, string(res))
	}
}

func addNewVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	newVideo := &defs.NewVideo{}
	if err := json.Unmarshal(res, newVideo); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	videoInfo, err := dbops.AddVideoInfo(newVideo.Name, newVideo.AuthorId)
	if err != nil {
		log.Printf("Error in AddNewVideo: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	if res, err := json.Marshal(videoInfo); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 201, string(res))
	}
}

func listAllVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	userName := ps.ByName("user_name")
	videos, err := dbops.ListVideoInfo(userName, 0, time.Now().Unix())
	if err != nil {
		log.Printf("Error in ListAllVideos: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	resp := &defs.VideosResp{Videos: videos}
	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 200, string(res))
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	videoId := ps.ByName("video_id")
	err := dbops.DeleteVideoInfo(videoId)
	if err != nil {
		log.Printf("Error in DeleteVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	dbops.RequestDeleteVideoFile(videoId)
	sendNormalResponse(w, 204, "")
}

func postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	newComment := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, newComment); err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	videoId := ps.ByName("video_id")
	if _, err := dbops.AddComment(newComment.Content, newComment.AuthorId, videoId); err != nil {
		log.Printf("Error in PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
	} else {
		sendNormalResponse(w, 201, "comment post success")
	}
}

func listComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	videoId := ps.ByName("video_id")
	comments, err := dbops.ListComments(videoId, 0, time.Now().Unix())
	if err != nil {
		log.Printf("Error in ShowComments: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	resp := &defs.CommentsResp{Comments: comments}
	if resp, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, defs.ErrorInternalError)
	} else {
		sendNormalResponse(w, 200, string(resp))
	}
}
