package createComment

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

type CreateComment struct {
}

func New() *CreateComment {
	return &CreateComment{}
}

var commentsCollection *mongo.Collection = api.GetCollection(api.DB, "comments")
var postCollection *mongo.Collection = api.GetCollection(api.DB, "posts")

func (e *CreateComment) Status(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var jsonComment models.Comment
	err := json.NewDecoder(c.Request().Body).Decode(&jsonComment)
	if err != nil {
		return err
	}

	newComment := models.Comment{
		Id:       primitive.NewObjectID(),
		PostId:   jsonComment.PostId,
		Author:   jsonComment.Author,
		Icon:     jsonComment.Icon,
		Date:     jsonComment.Date,
		Text:     jsonComment.Text,
		IsAnswer: jsonComment.IsAnswer,
		Likes:    0,
		Dislikes: 0,
		Answers:  []models.Comment{},
	}

	postId, _ := primitive.ObjectIDFromHex(jsonComment.PostId)

	if !newComment.IsAnswer {
		comment, err := commentsCollection.InsertOne(ctx, newComment)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
		resp := postCollection.FindOne(ctx, bson.M{"id": postId})
		var post models.Post
		resp.Decode(&post)
		fmt.Println(post)

		postCollection.UpdateOne(ctx, bson.M{"id": postId}, bson.M{"$set": bson.M{"comments": post.Comments + 1}})

		return c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: &echo.Map{
				"comment": comment,
			},
		})
	}

	commentCode := commentsCollection.FindOne(ctx, bson.M{"id": jsonComment.PostId})
	var comment models.Comment
	commentCode.Decode(&comment)

	newAnswers := comment.Answers
	newAnswers = append(newAnswers, newComment)
	commentsCollection.UpdateOne(ctx, bson.M{"id": jsonComment.PostId}, bson.M{"$set": bson.M{"answers": newAnswers}})

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"comment": comment,
		},
	})

}
