package commentDelete

import (
	"context"
	"encoding/json"
	"fmt"
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

type CommentDelete struct {
}

func New() *CommentDelete {
	return &CommentDelete{}
}

var commentsCollection *mongo.Collection = api.GetCollection(api.DB, "comments")
var postsCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *CommentDelete) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var jsonComment models.CommentCommunication
	err := json.NewDecoder(c.Request().Body).Decode(&jsonComment)
	if err != nil {
		return err
	}

	AId, _ := primitive.ObjectIDFromHex(jsonComment.AId)
	CId, _ := primitive.ObjectIDFromHex(jsonComment.CId)

	if jsonComment.AId == "" {
		var comment models.Comment
		commentsCollection.FindOne(ctx, bson.M{"id": CId}).Decode(&comment)
		id, _ := primitive.ObjectIDFromHex(comment.PostId)
		commentsCollection.DeleteOne(ctx, bson.M{"id": CId})

		var post models.Post
		postsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&post)
		postsCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"comments": post.Comments - 1}})
		fmt.Println(post.Comments)
		fmt.Println(comment)
	} else {
		var comment models.Comment
		commentsCollection.FindOne(ctx, bson.M{"id": CId}).Decode(&comment)
		for c := range comment.Answers {
			if comment.Answers[c].Id == AId {
				comment.Answers[c] = comment.Answers[len(comment.Answers)-1]
				comment.Answers = comment.Answers[:len(comment.Answers)-1]
			}
		}
		commentsCollection.UpdateOne(ctx, bson.M{"id": CId}, bson.M{"$set": bson.M{"answers": comment.Answers}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "commnet has been deleted",
		Data: &echo.Map{
			"status": "ok",
		},
	})
}
