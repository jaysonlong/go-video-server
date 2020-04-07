package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

// 解析api请求并修改转发到api服务器
func request(body *ApiBody, w http.ResponseWriter, r *http.Request) {
	var req *http.Request
	apiUrl := API_SERVER_HOST + body.Url

	switch body.Method {
	case "GET":
		req, _ = http.NewRequest("GET", apiUrl, nil)
	case "POST":
		req, _ = http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(body.RequestBody)))
	case "DELETE":
		req, _ = http.NewRequest("DELETE", apiUrl, nil)
	default:
		sendErrorResponse(w, ErrorRequestNotRecognized)
		return
	}

	req.Header = r.Header
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("%v", err)
		sendErrorResponse(w, ErrorInternalError)
		return
	}
	forwardResponse(w, resp)
}

func forwardResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func sendErrorResponse(w http.ResponseWriter, errResp ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	resp, _ := json.Marshal(&errResp.Err)
	io.WriteString(w, string(resp))
}
