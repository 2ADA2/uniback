package updateUser

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

type UpdateUser struct {
}

func New() *UpdateUser {
	return &UpdateUser{}
}

var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")
var cfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *UpdateUser) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var jsonUser models.UserCfg
	err := json.NewDecoder(c.Request().Body).Decode(&jsonUser)
	if err != nil {
		return err
	}

	authHeader := c.Request().Header.Get("Authorization")
	cursor := userCollection.FindOne(ctx, bson.M{"token": authHeader})
	var user models.User
	cursor.Decode(&user)

	if user.Name != jsonUser.User {
		return c.JSON(
			http.StatusOK, responses.UserResponse{
				Status:  http.StatusOK,
				Message: "wrong user",
			})
	}

	cfgCollection.UpdateOne(ctx, bson.M{"user": jsonUser.User}, bson.M{
		"$set": bson.M{
			"About": jsonUser.About,
			"Links": jsonUser.Links,
		},
	})

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"user": jsonUser,
		},
	})
}
