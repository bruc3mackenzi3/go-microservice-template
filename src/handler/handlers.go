package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func postUser(c echo.Context) error {
	return c.String(http.StatusOK, "POST user called")
}

func getUser(c echo.Context) error {
	return c.String(http.StatusOK, "GET user called")
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE user called")
}
