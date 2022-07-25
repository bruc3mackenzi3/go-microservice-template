package handler

import "github.com/labstack/echo"

func RegisterRoutes(e *echo.Echo) {
	e.POST("/users/:name", postUser)
	e.GET("/users/:id", getUser)
	e.DELETE("/users/:id", deleteUser)
}
