package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID
	Name     string
	Password string
	Token    string
}

type Post struct {
	ID        primitive.ObjectID
	Author    string `json:"author"`
	Subs      int    `json:"subs"`
	Date      string `json:"date"`
	Text      string `json:"text"`
	ImgUrl    string `json:"imgUrl"`
	Likes     int
	Bookmarks int
	Views     int
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Link struct {
	Name string
	Link string
}

type UserCfg struct {
	User       string
	About      string
	Followers  []string
	Subscribes []string
	Links      []Link
	Date       string
	Posts      Posts
	Likes      []string
	Bookmarks  []string
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
