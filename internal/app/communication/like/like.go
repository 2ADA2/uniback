package like

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

type Like struct {
}

func New() *Like {
	return &Like{}
}

var postCollection *mongo.Collection = api.GetCollection(api.DB, "posts")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")
var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func (e *Like) Status(c echo.Context) error {
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
			Message: "cannot like",
			Data: &echo.Map{
				"status": "post hasn't been found",
			},
		})
	}
	userLikes := cfg.Likes

	liked := false
	for i := range userLikes {
		if userLikes[i] == postId {
			liked = true
			userLikes[i] = userLikes[len(userLikes)-1]
			userLikes = userLikes[:len(userLikes)-1]
			break
		}
	}

	if !liked {
		userLikes = append(userLikes, postId)
		postCollection.UpdateOne(ctx, bson.M{"id": jsonPost.ID}, bson.M{"$set": bson.M{"likes": post.Likes + 1}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"likes": userLikes}})
	} else {
		postCollection.UpdateOne(ctx, bson.M{"id": jsonPost.ID}, bson.M{"$set": bson.M{"likes": post.Likes - 1}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"likes": userLikes}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "liked",
		Data: &echo.Map{
			"liked": liked,
		},
	})
}
