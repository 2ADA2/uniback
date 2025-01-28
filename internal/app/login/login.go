package login

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"myapp/internal/pkg/token"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Login struct {
}

func New() *Login {
	return &Login{}
}

var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func (e *Login) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name := c.FormValue("name")
	password := c.FormValue("password")

	if name == "" || password == "" {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": "Name and password are required"},
		})
	}

	userCode := userCollection.FindOne(ctx, bson.M{"name": name})

	if userCode == nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "not found",
			Data:    &echo.Map{"data": "not found"},
		})
	}

	var user models.User

	userCode.Decode(&user)

	if user.Password != password {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "incorrect",
			Data:    &echo.Map{"data": "incorrect"},
		})
	}

	jwtToken, err := token.CreateToken(name)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "token error",
			Data:    &echo.Map{"data": "token hasn't been created"},
		})
	}

<<<<<<< HEAD
	userCollection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"token": jwtToken}})
	resp := `{"Token" : ` + jwtToken + "}"
=======
	update := user
	update.Token = jwtToken

	userCollection.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": update})
	resp := `{"token" : ` + jwtToken + "}"
>>>>>>> 90a6f8c68d6526f03951040c77102c24ec1bb0df
	c.Response().Write([]byte(resp))

	return nil
}
