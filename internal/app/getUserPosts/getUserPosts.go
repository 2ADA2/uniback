package getuserposts

import (
	"context"
	"encoding/json"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetUserPosts struct {
}

func New() *GetUserPosts {
	return &GetUserPosts{}
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")
var usersCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func (e *GetUserPosts) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	var jsonUser models.UserInfo
	err := json.NewDecoder(c.Request().Body).Decode(&jsonUser)
	if err != nil {
		return err
	}
	var user models.User
	userCode := usersCollection.FindOne(ctx, bson.M{"name": jsonUser.User})
	userCode.Decode(&user)
	if user.Name == "" {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	var userPosts []models.Post
	postsCode, err := postsCollection.Find(ctx, bson.M{"author": user.Name})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	defer postsCode.Close(ctx)
	for postsCode.Next(ctx) {
		var userPost models.Post
		if err = postsCode.Decode(&userPost); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		userPosts = append(userPosts, userPost)
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"posts": userPosts,
		},
	})
}
