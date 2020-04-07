package dbops

import (
	"database/sql"
	"log"
	"time"

	"github.com/midmis/go-video-server/api/defs"
	"github.com/midmis/go-video-server/api/utils"
)

func AddUser(userName string, pwd string) error {
	stmt, err := dbConn.Prepare("INSERT INTO users (user_name, pwd) VALUES(?, ?)")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmt.Exec(userName, pwd)
	defer stmt.Close()
	if err != nil {
		log.Printf("%s %T", err, err)
		return err
	}

	return nil
}

func GetUserInfo(userName string) (*defs.UserInfo, error) {
	stmt, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE user_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var pwd string
	var userId int
	err = stmt.QueryRow(userName).Scan(&userId, &pwd)
	defer stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%s", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	userInfo := &defs.UserInfo{UserId: userId, UserName: userName, Pwd: pwd}
	return userInfo, nil
}

func DeleteUser(userName string, pwd string) error {
	stmt, err := dbConn.Prepare("DELETE FROM users WHERE user_name = ? AND pwd = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmt.Exec(userName, pwd)
	defer stmt.Close()
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	return nil
}

func AddVideoInfo(name string, authorId int) (*defs.VideoInfo, error) {
	ctime := time.Now().Format("Jan 02 2006, 15:04:05")
	videoId, _ := utils.NewUUID()

	stmt, err := dbConn.Prepare("INSERT INTO video_info (id, name, author_id, display_ctime) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	_, err = stmt.Exec(videoId, name, authorId, ctime)
	defer stmt.Close()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	videoInfo := &defs.VideoInfo{VideoId: videoId, Name: name, AuthorId: authorId, DisplayCtime: ctime}
	return videoInfo, nil
}

func GetVideoInfo(videoId string) (*defs.VideoInfo, error) {
	stmt, err := dbConn.Prepare("SELECT name, author_id, display_ctime FROM video_info WHERE id = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var (
		name     string
		authorId int
		ctime    string
	)
	err = stmt.QueryRow(videoId).Scan(&name, &authorId, &ctime)
	defer stmt.Close()

	if err != nil && err != sql.ErrNoRows {
		log.Printf("%s", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	videoInfo := &defs.VideoInfo{VideoId: videoId, Name: name, AuthorId: authorId, DisplayCtime: ctime}
	return videoInfo, nil
}

func ListVideoInfo(userName string, from, to int64) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT v.id, v.author_id, v.name, v.display_ctime FROM video_info v
		INNER JOIN users u ON v.author_id = u.id
		WHERE u.user_name=? AND v.create_time > FROM_UNIXTIME(?) AND v.create_time<=FROM_UNIXTIME(?)
		ORDER BY v.create_time DESC`)

	res := []*defs.VideoInfo{}
	if err != nil {
		return res, err
	}

	rows, err := stmtOut.Query(userName, from, to)
	defer stmtOut.Close()
	if err != nil {
		log.Printf("%v", err)
		return res, err
	}

	for rows.Next() {
		var videoId, name, ctime string
		var authorId int
		if err := rows.Scan(&videoId, &authorId, &name, &ctime); err != nil {
			return res, err
		}
		videoInfo := &defs.VideoInfo{VideoId: videoId, AuthorId: authorId, Name: name, DisplayCtime: ctime}
		res = append(res, videoInfo)
	}
	return res, nil
}

func DeleteVideoInfo(videoId string) error {
	stmt, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmt.Exec(videoId)
	defer stmt.Close()
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	return nil
}

func AddComment(content string, authorId int, videoId string) (*defs.Comment, error) {
	ctime := time.Now().Format("Jan 02 2006, 15:04:05")

	stmt, err := dbConn.Prepare("INSERT INTO comments (video_id, author_id, content, display_ctime) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	_, err = stmt.Exec(videoId, authorId, content, ctime)
	defer stmt.Close()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	comment := &defs.Comment{VideoId: videoId, Content: content, AuthorId: authorId, DisplayCtime: ctime}
	return comment, nil
}

func ListComments(videoId string, from, to int64) ([]*defs.Comment, error) {
	stmt, err := dbConn.Prepare(`SELECT c.id, c.author_id, u.user_name, c.content, c.display_ctime
		FROM comments c INNER JOIN users u ON c.author_id = u.id
		WHERE c.video_id = ? AND c.create_time > FROM_UNIXTIME(?) AND c.create_time <= FROM_UNIXTIME(?)
		ORDER BY c.create_time desc`)

	res := []*defs.Comment{}
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	rows, err := stmt.Query(videoId, from, to)
	defer stmt.Close()
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%s", err)
		return res, err
	}

	if err == sql.ErrNoRows {
		return res, nil
	}

	var (
		commentId  int
		authorId   int
		authorName string
		content    string
		ctime      string
	)
	for rows.Next() {
		rows.Scan(&commentId, &authorId, &authorName, &content, &ctime)
		comment := &defs.Comment{CommentId: commentId, VideoId: videoId,
			Content: content, AuthorId: authorId, AuthorName: authorName, DisplayCtime: ctime}
		res = append(res, comment)
	}

	return res, nil
}
