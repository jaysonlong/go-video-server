package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/midmis/go-video-server/scheduler/dbops"
)

func videoDeletionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	videoId := ps.ByName("video_id")

	if len(videoId) == 0 {
		sendResponse(w, http.StatusBadRequest, "video id cann't be empty")
		return
	}

	err := dbops.AddVideoDeletionRecord(videoId)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "internal service error")
		return
	}

	sendResponse(w, http.StatusOK, "delete action has been recorded")
}
