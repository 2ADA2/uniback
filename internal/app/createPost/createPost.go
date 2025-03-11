package createPost

import (
	"context"
	"encoding/json"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreatePost struct {
}

func New() *CreatePost {
	return &CreatePost{}
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *CreatePost) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	json_map := make(map[string]interface{})
	var jsonPost models.Post
	err := json.NewDecoder(c.Request().Body).Decode(&jsonPost)
	if err != nil {
		return err
	}

	newPost := models.Post{
		ID:        primitive.NewObjectID(),
		Author:    jsonPost.Author,
		Icon:      jsonPost.Icon,
		Header:    jsonPost.Header,
		Date:      jsonPost.Date,
		Text:      jsonPost.Text,
		ImgUrl:    jsonPost.ImgUrl,
		Likes:     jsonPost.Likes,
		Bookmarks: jsonPost.Bookmarks,
		Views:     jsonPost.Views,
	}

	p, err := postsCollection.InsertOne(ctx, newPost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"post": p,
			"id":   json_map["ID"],
		},
	})
}
