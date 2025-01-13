package ping

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Ping struct {
}

func New() *Ping {
	return &Ping{}
}

func (e *Ping) Status(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "pong")
	if err != nil {
		return err
	}

	return nil
}
