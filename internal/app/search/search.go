package search

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"myapp/internal/pkg/searchSort"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Search struct {
}

func New() *Search {
	return &Search{}
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")
var cfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *Search) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	search := c.QueryParam("searchValue")

	var searchResponse []interface{}

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

	for i := 0; i < len(posts); i++ {
		if strings.Contains(posts[i].Text, search) || strings.Contains(posts[i].Header, search) {
			searchResponse = append(searchResponse, posts[i])

		}
	}

	var users []models.UserCfg

	cursor, err = cfgCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cfg models.UserCfg
		if err := cursor.Decode(&cfg); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
		users = append(users, cfg)
	}

	for i := 0; i < len(users); i++ {
		if strings.Contains(users[i].User, search) {
			searchResponse = append(searchResponse, users[i])
		}
	}

	sortedSearch := searchSort.SearchSort(search, searchResponse)

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": sortedSearch},
	})
}
