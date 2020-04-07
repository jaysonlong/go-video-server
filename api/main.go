package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleware struct {
	router *httprouter.Router
}

func NewMiddleware(router *httprouter.Router) *middleware {
	return &middleware{router: router}
}

func (m middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ValidateSession(r)
	m.router.ServeHTTP(w, r)
}

func registerHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", createUser)
	router.POST("/user/:user_name", login)
	router.GET("/user/:user_name", getUserId)
	router.DELETE("/user", logout)
	router.POST("/user/:user_name/videos", addNewVideo)
	router.GET("/user/:user_name/videos", listAllVideos)
	router.DELETE("/user/:user_name/videos/:video_id", deleteVideo)
	router.POST("/videos/:video_id/comments", postComment)
	router.GET("/videos/:video_id/comments", listComments)

	return router
}

func main() {
	router := registerHandlers()
	middleWare := NewMiddleware(router)
	log.Fatal(http.ListenAndServe(":8000", middleWare))
}
