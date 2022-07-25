package handler

import (
	"net/http"

	"github.com/bruc3mackenzi3/microservice-demo/src/service"
	"github.com/labstack/echo"
)

var mService service.Service

func init() {
	mService = service.NewService()
}

func postUser(c echo.Context) error {
	mService.CreateUser("")
	return c.String(http.StatusOK, "POST user called")
}

func getUser(c echo.Context) error {
	mService.GetUser("")
	return c.String(http.StatusOK, "GET user called")
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE user called")
}
