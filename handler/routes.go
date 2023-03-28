package handler

import (
	"github.com/bruc3mackenzi3/microservice-demo/service"
	"github.com/labstack/echo/v4"
)

const usersPath = "/v1/users"

func RegisterRoutes(e *echo.Echo, s service.Service) {
	g := e.Group(usersPath)
	// TODO: add middleware to the group

	g.POST("/", postUser(s))
	g.GET("/:id", getUser(s))
	g.GET("/", getUsers(s))
	g.PUT("/:id", putUser(s))
	g.DELETE("/:id", deleteUser(s))
}
