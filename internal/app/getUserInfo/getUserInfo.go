package getUserInfo

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetUserInfo struct {
}

func New() *GetUserInfo {
	return &GetUserInfo{}
}

var usersCollection *mongo.Collection = api.GetCollection(api.DB, "users")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *GetUserInfo) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	name := c.QueryParam("name")

	var user models.User
	usersCollection.FindOne(ctx, bson.M{"name": name}).Decode(&user)

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
			"data": map[string]interface{}{
				"User":       userCfg.User,
				"About":      userCfg.About,
				"Followers":  userCfg.Followers,
				"Subscribes": userCfg.Subscribes,
				"Links":      userCfg.Links,
				"Date":       userCfg.Date,
			},
		},
	})
}
