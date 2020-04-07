package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/midmis/go-video-server/api/defs"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	resp, _ := json.Marshal(&errResp.Err)
	io.WriteString(w, string(resp))
}

func sendNormalResponse(w http.ResponseWriter, httpSC int, resp string) {
	w.WriteHeader(httpSC)
	io.WriteString(w, resp)
}
