package handler

import "github.com/labstack/echo"

const usersPath = "/users"

func RegisterRoutes(e *echo.Echo) {
	g := e.Group(usersPath)
	// TODO: add middleware to the group

	g.POST("/", postUser)
	g.GET("/:id", getUser)
	g.PUT("/:id", putUser)
	g.DELETE("/:id", deleteUser)
}
