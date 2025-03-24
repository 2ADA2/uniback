package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID
	Name     string
	Password string
	Token    string
}

type Post struct {
	ID        primitive.ObjectID
	Author    string `json:"author"`
	Icon      string
	Header    string
	Date      string `json:"date"`
	Text      string `json:"text"`
	ImgUrl    string `json:"imgUrl"`
	Likes     int
	Bookmarks int
	Reposts   int
	Views     int
	Comments  int
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Link struct {
	Name string
	Link string
}

type UserCfg struct {
	User            string
	Icon            string
	About           string
	Followers       []string
	Subscribes      []string
	Links           []Link
	Date            string
	Posts           Posts
	Likes           []string
	Reposts         []string
	Bookmarks       []string
	CommentLikes    []string
	CommentDislikes []string
}

type SubscribeBody struct {
	Author string
}

type ViewBody struct {
	Id string
}

type UserInfo struct {
	User string
}

type Comment struct {
	Id        primitive.ObjectID
	CommentId string
	PostId    string
	Author    string
	Icon      string
	Text      string
	Date      string
	Likes     int
	Dislikes  int
	IsAnswer  bool
	Answers   []Comment
}
type GetComments struct {
	Id string
}
type CommentCommunication struct {
	User string
	CId  string
	AId  string
}
