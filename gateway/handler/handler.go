package handler

import (
	"errors"
	"net/http"

	"github.com/airlangga-hub/food-delivery-app/gateway/helper"
	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type UserService interface {
	RegisterCustomer(user model.UserRegister) (model.UserInfo, error)
	Login(email, password string) (string, error)
	TopUpBalance(userID uuid.UUID, amount int) (model.PaymentLink, error)
	GetUserInfo(userID uuid.UUID) (model.UserInfo, error)
}

type OrderService interface {
	CreateOrder(userID uuid.UUID, items []model.Item) (model.Order, error)
	GetDrivers(orderID uuid.UUID) (model.FindDriver, error)
	ChooseDriver(orderID, driverID uuid.UUID) (model.Order, error)
	GetOrders(userID uuid.UUID) ([]model.Order, error)
	GiveRating(orderID uuid.UUID) error
	DriverGetPendingOrders() ([]model.Order, error)
	DriverApplyToTakeOrder(driverID, orderID uuid.UUID) error
	MarkOrderAsDone(orderID, driverID uuid.UUID) error
}

type Handler struct {
	UserSvc  UserService
	OrderSvc OrderService
	Validate *validator.Validate
}

func New(userSvc UserService, orderSvc OrderService, val *validator.Validate) *Handler {
	return &Handler{UserSvc: userSvc, OrderSvc: orderSvc, Validate: val}
}

func (h *Handler) RegisterCustomer(c *echo.Context) error {
	var payload RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	user := model.UserRegister{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
		Address:   payload.Address,
	}

	userInfo, err := h.UserSvc.RegisterCustomer(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "register failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    userInfo,
	})
}

func (h *Handler) Login(c *echo.Context) error {
	var payload LoginRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	token, err := h.UserSvc.Login(payload.Email, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password").Wrap(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) TopUpBalance(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload TopUpBalanceRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	paymentLink, err := h.UserSvc.TopUpBalance(claims.UserID, payload.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "top up failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    paymentLink,
	})
}

func (h *Handler) GetUserInfo(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	userInfo, err := h.UserSvc.GetUserInfo(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "get user info failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    userInfo,
	})
}

func (h *Handler) CreateOrder(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload CreateOrderRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	items := make([]model.Item, len(payload.Items))

	for i, item := range payload.Items {
		itemID, err := uuid.Parse(item.ItemID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid item id format").Wrap(err)
		}

		items[i] = model.Item{
			ID:       itemID,
			Quantity: item.Quantity,
		}
	}

	order, err := h.OrderSvc.CreateOrder(claims.UserID, items)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "create order failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    order,
	})
}

func (h *Handler) GetDrivers(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id format").Wrap(err)
	}

	findDrivers, err := h.OrderSvc.GetDrivers(orderID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no drivers found, we'll keep looking...").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "get drivers failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    findDrivers,
	})
}

func (h *Handler) ChooseDriver(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload ChooseDriverRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id format").Wrap(err)
	}

	driverID, err := uuid.Parse(payload.DriverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid driver id format").Wrap(err)
	}

	order, err := h.OrderSvc.ChooseDriver(orderID, driverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "choose driver failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    order,
	})
}

func (h *Handler) GetOrders(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	orders, err := h.OrderSvc.GetOrders(claims.UserID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "you haven't made any orders yet").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "get orders failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    orders,
	})
}

func (h *Handler) GiveRating(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload GiveRatingRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params").Wrap(err)
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id format").Wrap(err)
	}

	if err := h.OrderSvc.GiveRating(orderID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "give rating failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
	})
}

func (h *Handler) DriverGetPendingOrders(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	orders, err := h.OrderSvc.DriverGetPendingOrders()
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no pending orders found").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "get pending orders failed").Wrap(err)
	}

	return c.JSON(http.StatusOK, Response{
		Message: http.StatusText(http.StatusOK),
		Data:    orders,
	})
}

func (h *Handler) DriverApplyToTakeOrder(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id format").Wrap(err)
	}

	if err := h.OrderSvc.DriverApplyToTakeOrder(claims.UserID, orderID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "apply to take order failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
	})
}

func (h *Handler) DriverCompleteOrder(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != helper.RoleDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order id format").Wrap(err)
	}

	if err := h.OrderSvc.MarkOrderAsDone(orderID, claims.UserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "mark order as done failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
	})
}
