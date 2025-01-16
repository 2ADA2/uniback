package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID
	Name     string
	Password string
	Token    string
}

type Post struct {
	Author string `json:"author"`
	Subs   int    `json:"subs"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	ImgUrl string `json:"imgUrl"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}
