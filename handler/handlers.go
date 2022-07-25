package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bruc3mackenzi3/microservice-demo/model"
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
	ID string `json:"id"`
}

type getUserResponse struct {
	Name string `json:"name"`
}

func postUser(c echo.Context) error {
	name := c.Param("id")
	var response postUserResponse
	var err error

	response.ID, err = mService.CreateUser(name)
	if err != nil {
		r := errorResponse{500, "server error occured"}
		return c.JSON(r.Status, r)
	}
	return c.JSON(http.StatusOK, response)
}

func getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		fmt.Printf("Invalid id argument value='%s' supplied: %v", c.Param("id"), err)
		r := errorResponse{400, "id must be an unsigned integer"}
		return c.JSON(r.Status, r)
	}

	user, err := mService.GetUser(uint(id))
	if err != nil {
		var r errorResponse
		_, ok := err.(model.NotFoundError)
		if ok {
			r = errorResponse{404, "user not found"}
		} else {
			r = errorResponse{500, "server error occured"}
		}
		return c.JSON(r.Status, r)
	}

	response := getUserResponse{user.Name}

	return c.JSON(http.StatusOK, response)
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE user called")
}
