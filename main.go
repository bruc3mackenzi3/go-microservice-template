package main

import (
	"fmt"
	"net/http"

	"github.com/bruc3mackenzi3/microservice-demo/src/handler"
	"github.com/labstack/echo"
)

const port = 8080

func setupServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler.RegisterRoutes(e)

	err := e.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func main() {
	fmt.Printf("Starting web server on port %d...\n", port)

	setupServer()
}
