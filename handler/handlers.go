package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
	"github.com/bruc3mackenzi3/microservice-demo/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

var validate *validator.Validate
var mService service.Service

func init() {
	validate = validator.New()
	mService = service.NewService(repository.NewRepository())
}

type userRequestBody struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone"`
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func newUserResponseFromModel(u *model.User) userResponse {
	return userResponse{
		ID:    strconv.Itoa(int(u.ID)),
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func postUser(c echo.Context) error {
	var rUser userRequestBody

	// Parse request body to struct; will catch malformed JSON errors
	err := c.Bind(&rUser)
	if err != nil {
		c.Logger().Warn("Failed to decode request body: ", err)
		r := errorResponse{400, "bad request"}
		return c.JSON(r.Status, r)
	}

	// Validate struct based on validate tags in struct definition
	err = validate.Struct(rUser)
	if err != nil {
		c.Logger().Warn("Failed to validate request struct: ", err)
		r := errorResponse{400, "bad request"}
		return c.JSON(r.Status, r)
	}

	user := model.User{
		Name:  rUser.Name,
		Email: strings.ToLower(rUser.Email),
		Phone: rUser.Phone,
	}

	err = mService.CreateUser(&user)
	if err == model.ErrUserEmailTaken {
		c.Logger().Warnf("Cannot create user, email %s already taken", user.Email)
		r := errorResponse{400, "email already taken"}
		return c.JSON(r.Status, r)
	}
	if err != nil {
		c.Logger().Error("Failed to create user: ", err)
		r := errorResponse{500, "server error occured"}
		return c.JSON(r.Status, r)
	}

	response := newUserResponseFromModel(&user)
	return c.JSON(http.StatusOK, response)
}

func getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		c.Logger().Warn("Invalid id argument value='%s' supplied: %v", c.Param("id"), err)
		r := errorResponse{400, "id must be an unsigned integer"}
		return c.JSON(r.Status, r)
	}

	user, err := mService.GetUser(uint(id))
	if err != nil {
		var r errorResponse
		if err == model.ErrUserNotFound {
			r = errorResponse{404, "user not found"}
		} else if err != nil {
			r = errorResponse{500, "server error occured"}
		}
		return c.JSON(r.Status, r)
	}

	response := newUserResponseFromModel(user)
	return c.JSON(http.StatusOK, response)
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "DELETE user called")
}
