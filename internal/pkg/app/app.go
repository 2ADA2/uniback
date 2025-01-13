package app

import (
	"fmt"
	"log"
	"myapp/internal/app/controllers"
	getPosts "myapp/internal/app/endpoint"
	"myapp/internal/app/getUsers"
	"myapp/internal/app/ping"
	"myapp/internal/app/service"

	"github.com/labstack/echo/v4"
)

type App struct {
	e          *getPosts.GetPosts
	s          *service.Service
	ping       *ping.Ping
	getUsers   *getUsers.GetUsers
	createUser *controllers.CreateUser

	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.s = service.New()

	a.e = getPosts.New(a.s)
	a.ping = ping.New()
	a.getUsers = getUsers.New()
	a.createUser = controllers.New()

	a.echo = echo.New()

	a.echo.GET("/getPosts", a.e.Status)
	a.echo.GET("/ping", a.ping.Status)
	a.echo.POST("/createUser", a.createUser.Status)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(":3001")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
