package commentLike

import (
	"myapp/internal/app/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommentLike struct {
}

func New() *CommentLike {
	return &CommentLike{}
}

func (e *CommentLike) Status(c echo.Context) error {
	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &echo.Map{
			"comment": "",
		},
	})
}
