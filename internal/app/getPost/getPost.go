package getPost

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetPost struct {
}

func New() *GetPost {
	return &GetPost{}
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *GetPost) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postId, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	var post models.Post
	postsCollection.FindOne(ctx, bson.M{"id": postId}).Decode(&post)

	if post.Author == "" {
		return c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &echo.Map{"data": "no such data"},
		})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"post": post,
		},
	})
}
