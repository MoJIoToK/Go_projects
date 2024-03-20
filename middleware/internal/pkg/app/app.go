package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"middleware/internal/app/endpoint"
	"middleware/internal/app/mw"
	"middleware/internal/app/service"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	//Creating an Instance of service.
	a.s = service.New()

	//Creating an Instance of endpoint.
	a.e = endpoint.New(a.s)

	//Creating an Instance of echo.
	a.echo = echo.New()

	a.echo.Use(mw.CheckRole)
	//a.echo.Use(mw.CheckRole)

	//Description of the path and its handler
	a.echo.GET("/status", a.e.Status)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("Server running!")

	//Starting server
	err := a.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
