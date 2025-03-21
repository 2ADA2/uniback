package commentDislike

import (
	"myapp/internal/app/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommentDislike struct {
}

func New() *CommentDislike {
	return &CommentDislike{}
}

func (e *CommentDislike) Status(c echo.Context) error {
	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"comment": "",
		},
	})
}
