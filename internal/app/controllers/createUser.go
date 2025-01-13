package controllers

import (
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

type CreateUser struct {
}

func New() *CreateUser {
	return &CreateUser{}
}
func (e *CreateUser) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get values from the request
	name := c.FormValue("name")
	password := c.FormValue("password")

	// Validate inputs
	if name == "" || password == "" {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": "Name and password are required"},
		})
	}

	// Create a new user object
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     name,
		Password: password,
	}

	// Insert the new user into the database
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	// Respond with the created user ID
	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &echo.Map{"data": result},
	})
}
