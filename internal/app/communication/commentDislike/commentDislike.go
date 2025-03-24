package commentDislike

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

type CommentDislike struct {
}

func New() *CommentDislike {
	return &CommentDislike{}
}

var commentsCollection *mongo.Collection = api.GetCollection(api.DB, "comments")
var userCfgCollection *mongo.Collection = api.GetCollection(api.DB, "configs")

func (e *CommentDislike) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var jsonComment models.CommentCommunication
	err := json.NewDecoder(c.Request().Body).Decode(&jsonComment)
	if err != nil {
		return err
	}

	var userCfg models.UserCfg
	userCfgCollection.FindOne(ctx, bson.M{"user": jsonComment.User}).Decode(&userCfg)
	dislikes := userCfg.CommentDislikes

	AId, _ := primitive.ObjectIDFromHex(jsonComment.AId)
	CId, _ := primitive.ObjectIDFromHex(jsonComment.CId)

	if jsonComment.AId == "" {
		var comment models.Comment
		commentsCollection.FindOne(ctx, bson.M{"id": CId}).Decode(&comment)
		liked := false
		id := comment.Id.Hex()
		for c := range dislikes {
			if id == dislikes[c] {
				liked = true
				comment.Dislikes -= 1
				dislikes[c] = dislikes[len(dislikes)-1]
				dislikes = dislikes[:len(dislikes)-1]
				break
			}
		}
		if !liked {
			dislikes = append(dislikes, CId.Hex())
			comment.Dislikes += 1
		}
		commentsCollection.UpdateOne(ctx, bson.M{"id": CId}, bson.M{"$set": bson.M{"likes": comment.Likes}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": jsonComment.User}, bson.M{"$set": bson.M{"commentDislikes": dislikes}})

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
				for c := range dislikes {
					if id == dislikes[c] {
						liked = true
						comment.Answers[i].Likes -= 1
						dislikes[c] = dislikes[len(dislikes)-1]
						dislikes = dislikes[:len(dislikes)-1]
						break
					}
				}
			}
		}

		if !liked {
			dislikes = append(dislikes, AId.Hex())
			comment.Answers[ans].Dislikes += 1
		}
		commentsCollection.UpdateOne(ctx, bson.M{"id": CId}, bson.M{"$set": bson.M{"answers": comment.Answers}})
		userCfgCollection.UpdateOne(ctx, bson.M{"user": jsonComment.User}, bson.M{"$set": bson.M{"commentDislikes": dislikes}})
	}
	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "comment has been disliked",
		Data: &echo.Map{
			"status": "ok",
		},
	})
}
