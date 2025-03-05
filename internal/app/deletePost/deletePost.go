package deletePost

import (
	"context"
	"myapp/internal/app/responses"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeletePost struct {
}

func New() *DeletePost {
	return &DeletePost{}
}

var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *DeletePost) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postId, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	postsCollection.DeleteOne(ctx, bson.M{"id": postId})

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "deleted",
		Data: &echo.Map{
			"data": "deleted",
		},
	})
}
