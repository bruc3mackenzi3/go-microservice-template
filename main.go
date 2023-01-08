package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func initServer() *echo.Echo {
	e := echo.New()

	// Register Kubernetes Liveness Probe
	e.GET("/healthz", handler.LivenessProbe)

	// Register routes of User resource
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

	e.Logger.Infof("Starting web server on port %d...\n", config.Port)

	return e
}

func startServer(e *echo.Echo) {
	err := e.Start(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		e.Logger.Warnf("Echo server stopped and returned an error: %v", err)
	}
}

// catchSignal listens on a channel for interrupts, and gracefully shuts down the web server
// before terminating the application.  It listens for the following signals:
// - SIGINT, which is sent from the terminal by a user hitting CTRL + C
// - SIGTERM, which is sent by Docker when terminating the container
func catchSignal(e *echo.Echo) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	signal := <-signalChan
	e.Logger.Infof("Received signal %v, shutting down server...", signal)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error("Echo Shutdown encountered an error:", err)
	} else {
		e.Logger.Info("Server successfully shutdown gracefully!")
	}
}

func main() {
	e := initServer()

	// Start Echo web server in a separate Goroutine
	go startServer(e)

	// Block on catching shutdown signals in the main thread
	catchSignal(e)
}
