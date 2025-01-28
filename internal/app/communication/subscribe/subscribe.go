package subscribe

import (
	"myapp/internal/app/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Subscribe struct {
}

func New() *Subscribe {
	return &Subscribe{}
}

func (e *Subscribe) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "subscribed",
		Data: &echo.Map{
			"added": true,
		},
	})
}
