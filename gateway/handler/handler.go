package handler

import (
	"net/http"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type UserService interface {
	UserRegister(user model.UserRegister) (model.UserInfo, error)
	UserLogin(email, password string) (string, error)
}

type OrderService interface{}

type Handler struct {
	UserSvc  UserService
	OrderSvc OrderService
	Validate *validator.Validate
}

func New(userSvc UserService, orderSvc OrderService, val *validator.Validate) *Handler {
	return &Handler{UserSvc: userSvc, OrderSvc: orderSvc, Validate: val}
}

func (h *Handler) Register(c *echo.Context) error {
	var request RegisterRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	user := model.UserRegister{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		Address:   request.Address,
	}

	userInfo, err := h.UserSvc.UserRegister(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "register failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    userInfo,
	})
}

func (h *Handler) Login(c *echo.Context) error {
	var request LoginRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	token, err := h.UserSvc.UserLogin(request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password").Wrap(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
