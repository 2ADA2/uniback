package getComments

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetComments struct {
}

func New() *GetComments {
	return &GetComments{}
}

var commentsCollection *mongo.Collection = api.GetCollection(api.DB, "comments")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *GetComments) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.QueryParam("id")

	cursor, err := commentsCollection.Find(ctx, bson.M{"postid": id})
	var comments []models.Comment

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
		var user models.UserCfg
		userCfgCollection.FindOne(ctx, bson.M{"user": comment.Author}).Decode(&user)
		comment.Icon = user.Icon
		comments = append(comments, comment)
	}
	slices.Reverse(comments)

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"comments": comments,
		},
	})
}
