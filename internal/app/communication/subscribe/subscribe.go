package subscribe

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

type Subscribe struct {
}

func New() *Subscribe {
	return &Subscribe{}
}

var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")
var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func (e *Subscribe) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var body models.SubscribeBody
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return err
	}

	token := c.Request().Header.Get("Authorization")

	var user models.User
	var author models.User

	userCode := userCollection.FindOne(ctx, bson.M{"token": token})
	authorCode := userCollection.FindOne(ctx, bson.M{"name": body.Author})

	userCode.Decode(&user)
	authorCode.Decode(&author)

	var userCfg models.UserCfg
	var authorCfg models.UserCfg

	userCfgCode := userCfgCollection.FindOne(ctx, bson.M{"user": user.Name})
	authorCfgCode := userCfgCollection.FindOne(ctx, bson.M{"user": body.Author})

	userCfgCode.Decode(&userCfg)
	authorCfgCode.Decode(&authorCfg)

	if userCfg.User == "" || authorCfg.User == "" {
		return c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "cannot find users",
			Data: &echo.Map{
				"subscribed": false,
			},
		})
	}

	subscribed := false
	followed := false
	followers := authorCfg.Followers
	subscribes := userCfg.Subscribes

	subId := -1
	followId := -1

	for i := range authorCfg.Followers {
		if authorCfg.Followers[i] == user.Name {
			subscribed = true
			subId = i
		}
	}

	for i := range userCfg.Subscribes {
		if userCfg.Subscribes[i] == author.Name {
			followed = true
			followId = i
		}
	}

	if subscribed != followed {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "no match found",
		})
	}

	if subscribed {
		newFollowers := followers[:followId]
		newFollowers = append(newFollowers, followers[followId+1:]...)

		newSubscribes := subscribes[:subId]
		newSubscribes = append(newSubscribes, subscribes[subId+1:]...)

		userCfgCollection.UpdateOne(ctx, bson.M{"user": body.Author}, bson.M{"$set": bson.M{"followers": newFollowers}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"subscribes": newSubscribes}})
	} else {
		followers = append(followers, user.Name)
		subscribes = append(subscribes, body.Author)
		userCfgCollection.UpdateOne(ctx, bson.M{"user": body.Author}, bson.M{"$set": bson.M{"followers": followers}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": user.Name}, bson.M{"$set": bson.M{"subscribes": subscribes}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "subscribed",
		Data: &echo.Map{
			"subscribed": followed,
		},
	})
}
