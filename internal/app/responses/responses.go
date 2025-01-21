package responses

import (
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}

type PostResponse struct {
	ID        string
	Author    string `json:"author"`
	Subs      int    `json:"subs"`
	Date      string `json:"date"`
	Text      string `json:"text"`
	ImgUrl    string `json:"imgUrl"`
	Likes     int
	Bookmarks int
	Views     int
}
