package main

// config
var API_SERVER_HOST = "http://127.0.0.1:8000"
var STREAM_SERVER_HOST = "http://127.0.0.1:9000"

// models
type UserPage struct {
	Name string
}

type ApiBody struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	RequestBody string `json:"request_body"`
}

type ErrorResponse struct {
	HttpSC int
	Err    Error
}

type Error struct {
	ErrMsg  string
	ErrCode string
}

var (
	ErrorRequestNotRecognized = ErrorResponse{
		HttpSC: 400,
		Err:    Error{"Api not recognize, bad request", "002"},
	}
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Err:    Error{"Request body not correct", "001"},
	}
	ErrorInternalError = ErrorResponse{
		HttpSC: 500,
		Err:    Error{"Internal service error", "004"},
	}
)
