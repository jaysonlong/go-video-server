package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/midmis/go-video-server/scheduler/taskrunner"
)

func registerHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:video_id", videoDeletionHandler)

	return router
}

func main() {
	go taskrunner.StartWorker()
	router := registerHandlers()
	log.Fatal(http.ListenAndServe(":9001", router))
}
