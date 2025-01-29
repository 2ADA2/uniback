package bookmark

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

type Bookmark struct {
}

func New() *Bookmark {
	return &Bookmark{}
}

var postCollection *mongo.Collection = api.GetCollection(api.DB, "posts")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")
var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func (e *Bookmark) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Request().Header.Get("Authorization")

	var jsonPost models.Post
	err := json.NewDecoder(c.Request().Body).Decode(&jsonPost)
	if err != nil {
		return err
	}

	resp := userCollection.FindOne(ctx, bson.M{"token": token})
	var user models.User
	resp.Decode(&user)

	resp = userCfgCollection.FindOne(ctx, bson.M{"user": user.Name})
	var cfg models.UserCfg
	resp.Decode(&cfg)

	postId := jsonPost.ID.Hex()
	resp = postCollection.FindOne(ctx, bson.M{"id": jsonPost.ID})
	var post models.Post
	resp.Decode(&post)

	if post.Author == "" {
		return c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "cannot add in bookmarks",
			Data: &echo.Map{
				"status": "post hasn't been found",
			},
		})
	}

	userBookmarks := cfg.Bookmarks

	added := false
	for i := range userBookmarks {
		if userBookmarks[i] == postId {
			added = true
			newBookmarks := userBookmarks[:i]
			newBookmarks = append(newBookmarks, userBookmarks[i+1:]...)
			userBookmarks = newBookmarks
			break
		}
	}

	if !added {
		userBookmarks = append(userBookmarks, postId)
		postCollection.UpdateOne(ctx, bson.M{"id": jsonPost.ID}, bson.M{"$set": bson.M{"bookmarks": post.Bookmarks + 1}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"bookmarks": userBookmarks}})
	} else {
		postCollection.UpdateOne(ctx, bson.M{"id": jsonPost.ID}, bson.M{"$set": bson.M{"bookmarks": post.Bookmarks - 1}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"bookmarks": userBookmarks}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "added",
		Data: &echo.Map{
			"added": added,
		},
	})
}
