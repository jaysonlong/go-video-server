package main

import (
	"io"
	"net/http"
)

func sendResponse(w http.ResponseWriter, HttpSC int, resp string) {
	w.WriteHeader(HttpSC)
	io.WriteString(w, resp)
}
