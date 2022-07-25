package handler

import (
	"net/http"

	"github.com/bruc3mackenzi3/microservice-demo/service"
	"github.com/labstack/echo"
)

var mService service.Service

func init() {
	mService = service.NewService()
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type postUserResponse struct {
	Id string `json:"name"`
}

func postUser(c echo.Context) error {
	name := c.Param("name")
	var response postUserResponse
	var err error

	response.Id, err = mService.CreateUser(name)
	if err != nil {
		r := errorResponse{500, "server error occured"}
		return c.JSON(r.Status, r)
	}
	return c.JSON(http.StatusOK, response)
}

func getUser(c echo.Context) error {
	mService.GetUser("")
	return c.String(http.StatusOK, "GET user called")
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE user called")
}
