package commentLike

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentLike struct {
}

func New() *CommentLike {
	return &CommentLike{}
}

var commentsCollection *mongo.Collection = api.GetCollection(api.DB, "comments")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *CommentLike) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var jsonComment models.CommentCommunication
	err := json.NewDecoder(c.Request().Body).Decode(&jsonComment)
	if err != nil {
		return err
	}

	var userCfg models.UserCfg
	userCfgCollection.FindOne(ctx, bson.M{"user": jsonComment.User}).Decode(&userCfg)
	likes := userCfg.CommentLikes

	AId, _ := primitive.ObjectIDFromHex(jsonComment.AId)
	CId, _ := primitive.ObjectIDFromHex(jsonComment.CId)

	if jsonComment.AId == "" {
		var comment models.Comment
		commentsCollection.FindOne(ctx, bson.M{"id": CId}).Decode(&comment)
		liked := false
		id := comment.Id.Hex()
		for c := range likes {
			if id == likes[c] {
				liked = true
				comment.Likes -= 1
				likes[c] = likes[len(likes)-1]
				likes = likes[:len(likes)-1]
				break
			}
		}
		if !liked {
			likes = append(likes, CId.Hex())
			comment.Likes += 1
		}
		commentsCollection.UpdateOne(ctx, bson.M{"id": CId}, bson.M{"$set": bson.M{"likes": comment.Likes}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": jsonComment.User}, bson.M{"$set": bson.M{"commentLikes": likes}})

	} else {
		var comment models.Comment
		commentsCollection.FindOne(ctx, bson.M{"id": CId}).Decode(&comment)
		liked := false
		var ans int
		for i := range comment.Answers {
			if comment.Answers[i].Id == AId {
				ans = i
				answer := comment.Answers[i]
				id := answer.Id.Hex()
				for c := range likes {
					if id == likes[c] {
						liked = true
						comment.Answers[i].Likes -= 1
						likes[c] = likes[len(likes)-1]
						likes = likes[:len(likes)-1]
						break
					}
				}
			}
		}

		if !liked {
			likes = append(likes, AId.Hex())
			comment.Answers[ans].Likes += 1
		}
		commentsCollection.UpdateOne(ctx, bson.M{"id": CId}, bson.M{"$set": bson.M{"answers": comment.Answers}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": jsonComment.User}, bson.M{"$set": bson.M{"commentLikes": likes}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "comment has been liked",
		Data: &echo.Map{
			"status": "ok",
		},
	})
}
