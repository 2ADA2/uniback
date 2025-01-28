package bookmark

import (
	"myapp/internal/app/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Bookmark struct {
}

func New() *Bookmark {
	return &Bookmark{}
}

func (e *Bookmark) Status(c echo.Context) error {
	return c.JSON(http.StatusAccepted, responses.UserResponse{
		Status:  http.StatusAccepted,
		Message: "added",
		Data: &echo.Map{
			"added": true,
		},
	})
}
