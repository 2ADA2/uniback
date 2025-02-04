package getUser

import (
	"context"
	"fmt"
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

type GetUser struct {
}

func New() *GetUser {
	return &GetUser{}
}

var usersCollection *mongo.Collection = api.GetCollection(api.DB, "users")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *GetUser) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Request().Header.Get("Authorization")
	userId, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	fmt.Println(userId)
	if c.QueryParam("id") != "" {
		fmt.Println(1)
	}

	var user models.User
	usersCollection.FindOne(ctx, bson.M{"token": token}).Decode(&user)

	if user.Name == "" {
		return c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "notFound",
			Data:    &echo.Map{"data": "no such data"},
		})
	}

	var userCfg models.UserCfg
	userCfgCollection.FindOne(ctx, bson.M{"user": user.Name}).Decode(&userCfg)

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"data": userCfg,
		},
	})
}
