package main

import (
	"fmt"
	"net/http"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func setupServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler.RegisterRoutes(e)

	e.Logger.SetLevel(log.DEBUG)

	err := e.Start(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func main() {
	fmt.Printf("Starting web server on port %d...\n", config.Port)

	setupServer()
}
