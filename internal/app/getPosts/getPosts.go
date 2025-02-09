package getPosts

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/app/service"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	GenerateNewPosts() service.Posts
}

type GetPosts struct {
	s Service
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func New(s Service) *GetPosts {
	return &GetPosts{
		s: s,
	}
}

func (e *GetPosts) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var posts []models.Post

	cursor, err := postsCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": posts},
	})
}
