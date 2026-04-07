package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, order model.Order) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
}

type PaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type userService struct {
	customerRepo       CustomerRepository
	paymentGatewayRepo PaymentGatewayRepository
}

func NewUserService(customerRepo CustomerRepository, paymentGatewayRepo PaymentGatewayRepository) *userService {
	return &userService{customerRepo: customerRepo, paymentGatewayRepo: paymentGatewayRepo}
}

func (s *userService) CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, order model.Order) (model.Order, model.PaymentGatewayResponse, error) {
	oorder, err := s.customerRepo.CreateOrder(ctx, userID, order)
	if err != nil {
		return model.Order{}, model.PaymentGatewayResponse{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
	}

	items := make([]model.PaymentGatewayItem, 0, 8)
	for _, resto := range oorder.Restaurants {
		for _, item := range resto.Items {
			items = append(items, model.PaymentGatewayItem{
				ReferenceID:   item.ID.String(),
				Name:          item.Name,
				Type:          "PHYSICAL_PRODUCT",
				Category:      "Order",
				NetUnitAmount: item.Price,
				Quantity:      item.Quantity,
			})
		}
	}

	paymentGatewayResponse, err := s.paymentGatewayRepo.CreatePaymentSession(ctx, userID, userEmail, order.TotalFee, items)
	if err != nil {
		return model.Order{}, model.PaymentGatewayResponse{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
	}

	oorder.PaymentLink = paymentGatewayResponse.PaymentLinkURL

	return oorder, paymentGatewayResponse, nil
}

func (s *userService) GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error) {
	drivers, err := s.customerRepo.GetDrivers(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order.customer_service.GetDrivers: %w", err)
	}
	return drivers, nil
}
