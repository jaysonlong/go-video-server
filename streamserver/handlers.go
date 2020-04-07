package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	videoId := ps.ByName("video_id")
	if len(videoId) == 0 {
		sendResponse(w, 400, "incorrect request content")
		return
	}

	videoPath := VIDEOS_DIR + videoId
	file, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Open video file error: %v", err)
		sendResponse(w, http.StatusInternalServerError, "internal service error")
		return
	}

	defer file.Close()
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), file)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendResponse(w, 400, "request body parse failed")
		return
	}

	srcFile, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Read request body file error: %v", err)
		sendResponse(w, http.StatusInternalServerError, "internal service error")
	}

	videoId := ps.ByName("video_id")
	if len(videoId) == 0 {
		sendResponse(w, 400, "incorrect request content")
		return
	}

	videoPath := VIDEOS_DIR + videoId
	dstFile, err := os.Create(videoPath)
	if err != nil {
		log.Printf("Create video file error: %v", err)
		sendResponse(w, http.StatusInternalServerError, "internal service error")
		return
	}

	io.Copy(dstFile, srcFile)
	sendResponse(w, http.StatusOK, "upload success")
}

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./upload.html")
	t.Execute(w, nil)
}
