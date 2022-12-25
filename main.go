package main

import (
	"fmt"
	"net/http"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func setupServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler.RegisterRoutes(e)

	// Enabling the middleware logger makes Echo log each http request received
	// e.Use(middleware.Logger())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

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
