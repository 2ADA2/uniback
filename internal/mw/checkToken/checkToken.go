package checkToken

import (
	"context"
	"myapp/internal/app/models"
	"myapp/internal/pkg/api"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = api.GetCollection(api.DB, "users")

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing token",
			})
		}

		cursor := userCollection.FindOne(ctx, bson.M{"token": authHeader})
		var user models.User
		cursor.Decode(&user)

		if user.Name == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "uncorrect token")
		}

		return next(c)
	}
}
