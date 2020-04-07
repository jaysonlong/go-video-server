package defs

// config
var VIDEO_DEL_SERVER_HOST = "http://127.0.0.1:9001"

// request body
type UserCredential struct {
	UserName string `json:"user_name"`
	Pwd      string `json:"pwd"`
}
type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}
type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

//response
type SignUpResp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}
type LoginResp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}
type LogoutResp struct {
	Success bool `json:"success"`
}
type VideosResp struct {
	Videos []*VideoInfo `json:"videos"`
}
type CommentsResp struct {
	Comments []*Comment `json:"comments"`
}
type UserInfoResp struct {
	UserId   int    `json:"id"`
	UserName string `json:"user_name"`
}

// data & response models
type VideoInfo struct {
	VideoId      string `json:"id"`
	Name         string `json:"name"`
	AuthorId     int    `json:"author_id"`
	DisplayCtime string `json:"display_ctime"`
}
type Comment struct {
	CommentId    int    `json:"id"`
	AuthorId     int    `json:"author_id"`
	AuthorName   string `json:"author"`
	VideoId      string `json:"video_id"`
	Content      string `json:"content"`
	DisplayCtime string `json:"display_ctime"`
}

// data models
type UserInfo struct {
	UserId   int
	UserName string
	Pwd      string
}

type SimpleSession struct {
	SessionId string
	LoginName string
	TTL       int
}
