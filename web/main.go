package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func registerHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)

	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)

	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler)

	router.GET("/videos/:video_id", proxyHandler)

	router.POST("/upload/:video_id", proxyHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	router := registerHandlers()
	log.Fatal(http.ListenAndServe(":8080", router))
}
