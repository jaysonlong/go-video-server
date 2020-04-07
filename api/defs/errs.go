package defs

type ErrorResponse struct {
	HttpSC int
	Err    Error
}

type Error struct {
	ErrMsg  string
	ErrCode string
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Err:    Error{"Request body not correct", "001"},
	}
	ErrorNotAuthUser = ErrorResponse{
		HttpSC: 401,
		Err:    Error{"User authentication failed", "002"},
	}
	ErrorDBError = ErrorResponse{
		HttpSC: 500,
		Err:    Error{"DB operation failed", "003"},
	}
	ErrorInternalError = ErrorResponse{
		HttpSC: 500,
		Err:    Error{"Internal service error", "004"},
	}
)
