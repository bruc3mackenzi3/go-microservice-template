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
	return c.JSON(http.StatusCreated, response)
}

func getUser(c echo.Context) error {
	userID, errResponse := parseID(c.Param("id"))
	if errResponse != nil {
		c.Logger().Warn("Invalid id argument value='%s'", c.Param("id"))
		return c.JSON(errResponse.Status, errResponse)
	}

	user, err := mService.GetUser(userID)
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

func putUser(c echo.Context) error {
	userID, errResponse := parseID(c.Param("id"))
	if errResponse != nil {
		c.Logger().Warn("Invalid id argument value='%s'", c.Param("id"))
		return c.JSON(errResponse.Status, errResponse)
	}

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
	user.ID = userID

	err = mService.UpdateUser(&user)
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
	return c.JSON(http.StatusCreated, response)
}

func deleteUser(c echo.Context) error {
	userID, errResponse := parseID(c.Param("id"))
	if errResponse != nil {
		c.Logger().Warn("Invalid id argument value='%s'", c.Param("id"))
		return c.JSON(errResponse.Status, errResponse)
	}

	err := mService.DeleteUser(userID)
	if err != nil {
		var r errorResponse
		if err == model.ErrUserNotFound {
			r = errorResponse{404, "user not found"}
		} else if err != nil {
			r = errorResponse{500, "server error occured"}
		}
		return c.JSON(r.Status, r)
	}

	return c.NoContent(http.StatusOK)
}

func parseID(idParam string) (uint, *errorResponse) {
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		return 0, &errorResponse{400, "id must be an unsigned integer"}
	}
	return uint(id), nil
}
