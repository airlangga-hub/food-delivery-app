package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/airlangga-hub/food-delivery-app/gateway/helper"
	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type UserService interface {
	RegisterCustomer(ctx context.Context, user model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email, password string) (string, error)
	TopUpBalance(ctx context.Context, userID string, amount int) (model.PaymentLink, error)
	GetUserInfo(ctx context.Context, userID string) (model.UserInfo, error)
	PaymentGatewayWebhook(ctx context.Context, userID string, paymentType model.PaymentType, amount int) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, userEmail string, deliveryFee int, items []model.Item) (model.Order, error)
	GetDrivers(ctx context.Context, orderID string) (model.FindDriver, error)
	ChooseDriver(ctx context.Context, orderID, driverID string) (model.Order, error)
	CustomerGetOrders(ctx context.Context, userID string) ([]model.Order, error)
	GiveRating(ctx context.Context, orderID string, rating int) error
	DriverGetPendingOrders(ctx context.Context) ([]model.Order, error)
	DriverApplyToTakeOrder(ctx context.Context, driverID, orderID string) error
	DriverCompleteOrder(ctx context.Context, orderID, driverID string) error
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

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	userInfo, err := h.UserSvc.RegisterCustomer(ctx, user)
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

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	token, err := h.UserSvc.Login(ctx, payload.Email, payload.Password)
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

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload TopUpBalanceRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	paymentLink, err := h.UserSvc.TopUpBalance(ctx, claims.UserID, payload.Amount)
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

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	userInfo, err := h.UserSvc.GetUserInfo(ctx, claims.UserID)
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

	if claims.Role != model.RoleUserCustomer {
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
		items[i] = model.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	order, err := h.OrderSvc.CreateOrder(ctx, claims.UserID, claims.Subject, payload.DeliveryFee, items)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "nothing found").Wrap(err)
		}
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

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order id must not be empty")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	findDrivers, err := h.OrderSvc.GetDrivers(ctx, orderID)
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

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload ChooseDriverRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order id must not be empty")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	order, err := h.OrderSvc.ChooseDriver(ctx, orderID, payload.DriverID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "nothing found").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "choose driver failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
		Data:    order,
	})
}

func (h *Handler) CustomerGetOrders(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	claims, ok := token.Claims.(*helper.MyClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	orders, err := h.OrderSvc.CustomerGetOrders(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no orders found").Wrap(err)
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

	if claims.Role != model.RoleUserCustomer {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user")
	}

	var payload GiveRatingRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params").Wrap(err)
	}

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order id must not be empty")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	if err := h.OrderSvc.GiveRating(ctx, orderID, payload.Rating); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "nothing found").Wrap(err)
		}
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

	if claims.Role != model.RoleUserDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	orders, err := h.OrderSvc.DriverGetPendingOrders(ctx)
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

	if claims.Role != model.RoleUserDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order id must not be empty")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	if err := h.OrderSvc.DriverApplyToTakeOrder(ctx, claims.UserID, orderID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "nothing found").Wrap(err)
		}
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

	if claims.Role != model.RoleUserDriver {
		return echo.NewHTTPError(http.StatusUnauthorized, "you're not a driver")
	}

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order id must not be empty")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	if err := h.OrderSvc.DriverCompleteOrder(ctx, orderID, claims.UserID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "nothing found").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "mark order as done failed").Wrap(err)
	}

	return c.JSON(http.StatusCreated, Response{
		Message: http.StatusText(http.StatusCreated),
	})
}

func (h *Handler) XenditWebhook(c *echo.Context) error {
	var payload XenditWebhookRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	if err := h.Validate.Struct(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").Wrap(err)
	}

	split := strings.Split(payload.Data.ReferenceID, "_")

	var paymentType model.PaymentType
	if strings.HasPrefix(payload.Data.ReferenceID, string(model.PaymentGatewayRefIDPrefixTopUp)) {
		paymentType = model.PaymentTypeTopUp
	} else if strings.HasPrefix(payload.Data.ReferenceID, string(model.PaymentGatewayRefIDPrefixOrder)) {
		paymentType = model.PaymentTypeOrder
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	if err := h.UserSvc.PaymentGatewayWebhook(ctx, split[1], paymentType, int(payload.Data.Amount)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "webhook error").Wrap(err)
	}
	
	return c.JSON(http.StatusOK, Response{Message: http.StatusText(http.StatusOK)})
}
