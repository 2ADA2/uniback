package getUsers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetUsers struct {
}

func New() *GetUsers {
	return &GetUsers{}
}

func (e *GetUsers) Status(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "pong")
	if err != nil {
		return err
	}

	return nil
}
