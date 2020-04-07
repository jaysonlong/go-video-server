package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleware struct {
	router  *httprouter.Router
	limiter *ConnLimiter
}

func NewMiddleware(router *httprouter.Router, concurrentCap int) *middleware {
	limiter := NewLimiter(concurrentCap)
	return &middleware{router: router, limiter: limiter}
}

func (m middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ValidateSession(r)
	if !m.limiter.openConn() {
		sendResponse(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	defer m.limiter.releaseConn()
	m.router.ServeHTTP(w, r)
}

func registerHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:video_id", streamHandler)
	router.POST("/upload/:video_id", uploadHandler)
	router.GET("/testpage", testPageHandler)

	return router
}

func main() {
	router := registerHandlers()
	middleWare := NewMiddleware(router, MAX_CONCURRENT_CNT)
	log.Fatal(http.ListenAndServe(":9000", middleWare))
}
