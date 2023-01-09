package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
	"github.com/bruc3mackenzi3/microservice-demo/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate
var userService service.Service

func init() {
	validate = validator.New()
	userService = service.NewService(repository.NewRepository())
}

type UserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone" validate:"omitempty,e164"`
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
	var rUser UserRequest

	// Parse request body to struct; will catch malformed JSON errors
	err := c.Bind(&rUser)
	if err != nil {
		c.Logger().Warn("Failed to decode request body: ", err)
		r := errorResponse{http.StatusBadRequest, "bad request"}
		return c.JSON(r.Status, r)
	}

	// Validate struct based on validate tags in struct definition
	err = validate.Struct(rUser)
	if err != nil {
		c.Logger().Warn("Failed to validate request struct: ", err)
		r := errorResponse{http.StatusBadRequest, "bad request"}
		return c.JSON(r.Status, r)
	}

	user := newUserFromRequest(rUser, 0)

	err = userService.CreateUser(&user)
	if err == model.ErrUserEmailTaken {
		c.Logger().Warnf("Cannot create user, email %s already taken", user.Email)
		r := errorResponse{http.StatusBadRequest, "email already taken"}
		return c.JSON(r.Status, r)
	} else if err != nil {
		c.Logger().Error("Failed to create user: ", err)
		r := errorResponse{http.StatusInternalServerError, "server error occured"}
		return c.JSON(r.Status, r)
	}

	response := newUserResponseFromModel(&user)
	return c.JSON(http.StatusCreated, response)
}

func getUser(c echo.Context) error {
	var userID uint
	err := echo.PathParamsBinder(c).Uint("id", &userID).BindError()
	if err != nil {
		c.Logger().Warn("Invalid id param:", err)
		return c.JSON(http.StatusBadRequest, errorResponse{http.StatusBadRequest, "bad user ID"})
	}

	user, err := userService.GetUser(userID)
	if err != nil {
		var r errorResponse
		if err == model.ErrUserNotFound {
			c.Logger().Warnf("User %s not found: %v", userID, err)
			r = errorResponse{http.StatusNotFound, "user not found"}
		} else if err != nil {
			c.Logger().Error("Failed to get User:", err)
			r = errorResponse{http.StatusInternalServerError, "server error occured"}
		}
		return c.JSON(r.Status, r)
	}
	c.Logger().Infof("Retrieved User %+v", user)

	response := newUserResponseFromModel(user)
	return c.JSON(http.StatusOK, response)
}

func putUser(c echo.Context) error {
	var userID uint
	err := echo.PathParamsBinder(c).Uint("id", &userID).BindError()
	if err != nil {
		c.Logger().Warn("Invalid id param:", err)
		return c.JSON(http.StatusBadRequest, errorResponse{http.StatusBadRequest, "bad user ID"})
	}

	var rUser UserRequest

	// Parse request body to struct; will catch malformed JSON errors
	err = c.Bind(&rUser)
	if err != nil {
		c.Logger().Warn("Failed to decode request body: ", err)
		r := errorResponse{http.StatusBadRequest, "bad request"}
		return c.JSON(r.Status, r)
	}

	// Validate struct based on validate tags in struct definition
	err = validate.Struct(rUser)
	if err != nil {
		c.Logger().Warn("Failed to validate request struct: ", err)
		r := errorResponse{http.StatusBadRequest, "bad request"}
		return c.JSON(r.Status, r)
	}

	user := newUserFromRequest(rUser, userID)

	err = userService.UpdateUser(&user)
	if err == model.ErrUserNotFound {
		c.Logger().Warnf("User %s not found: %v", userID, err)
		r := errorResponse{http.StatusNotFound, "user not found"}
		return c.JSON(r.Status, r)
	} else if err == model.ErrUserEmailTaken {
		c.Logger().Warnf("Cannot update user, email %s already taken", user.Email)
		r := errorResponse{http.StatusBadRequest, "email already taken"}
		return c.JSON(r.Status, r)
	}
	if err != nil {
		c.Logger().Error("Failed to update user: ", err)
		r := errorResponse{http.StatusInternalServerError, "server error occured"}
		return c.JSON(r.Status, r)
	}

	response := newUserResponseFromModel(&user)
	return c.JSON(http.StatusOK, response)
}

func deleteUser(c echo.Context) error {
	var userID uint
	err := echo.PathParamsBinder(c).Uint("id", &userID).BindError()
	if err != nil {
		c.Logger().Warn("Invalid id param:", err)
		return c.JSON(http.StatusBadRequest, errorResponse{http.StatusBadRequest, "bad user ID"})
	}

	err = userService.DeleteUser(userID)
	if err != nil {
		var r errorResponse
		if err == model.ErrUserNotFound {
			c.Logger().Warnf("User %s not found: %v", userID, err)
			r = errorResponse{http.StatusNotFound, "user not found"}
		} else if err != nil {
			c.Logger().Error("Failed to delete User: ", err)
			r = errorResponse{http.StatusInternalServerError, "server error occured"}
		}
		return c.JSON(r.Status, r)
	}

	return c.NoContent(http.StatusOK)
}

func newUserFromRequest(rUser UserRequest, userID uint) model.User {
	u := model.User{
		Name:  rUser.Name,
		Email: strings.ToLower(rUser.Email),
		Phone: rUser.Phone,
	}
	u.ID = userID
	return u
}
