package dbops

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "www:www@/go_video_server")

	if err != nil || dbConn.Ping() != nil {
		panic(err.Error())
	}
}
