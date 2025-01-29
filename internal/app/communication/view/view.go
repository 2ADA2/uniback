package view

import (
	"context"
	"encoding/json"
	"myapp/internal/app/models"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type View struct {
}

func New() *View {
	return &View{}
}

var postCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *View) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var body models.ViewBody
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return err
	}
	postId, _ := primitive.ObjectIDFromHex(body.Id)

	var post models.Post
	postCode := postCollection.FindOne(ctx, bson.M{"id": postId})
	postCode.Decode(&post)

	postCollection.UpdateOne(ctx, bson.M{"id": postId}, bson.M{"$set": bson.M{"views": post.Views + 1}})

	return c.String(http.StatusOK, "viewed")
}
