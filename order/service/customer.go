package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type SQLRepository interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, order model.OrderIn) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
	UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error
}

type PaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type MongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type customerService struct {
	sqlRepo            SQLRepository
	paymentGatewayRepo PaymentGatewayRepository
	mongoRepo          MongoRepository
}

func NewCustomerService(customerRepo SQLRepository, paymentGatewayRepo PaymentGatewayRepository, mongoRepo MongoRepository) *customerService {
	return &customerService{sqlRepo: customerRepo, paymentGatewayRepo: paymentGatewayRepo, mongoRepo: mongoRepo}
}

func (s *customerService) CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, order model.OrderIn) (model.Order, error) {
	oorder, err := s.sqlRepo.CreateOrder(ctx, userID, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
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

	paymentGatewayResponse, err := s.paymentGatewayRepo.CreatePaymentSession(ctx, model.PaymentTypeOrder, userID, userEmail, oorder.TotalFee, items)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
	}

	if err := s.mongoRepo.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailStatus:            model.EmailStatusPending,
		PaymentType:            model.PaymentTypeOrder,
		PaymentGatewayResponse: paymentGatewayResponse,
	}); err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
	}

	oorder.PaymentLink = paymentGatewayResponse.PaymentLinkURL

	return oorder, nil
}

func (s *customerService) GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error) {
	drivers, err := s.sqlRepo.GetDrivers(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order.customer_service.GetDrivers: %w", err)
	}
	return drivers, nil
}

func (s *customerService) PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error {
	var reason model.LedgerReason
	switch paymentType {
	case model.PaymentTypeOrder:
		reason = model.LedgerReasonCustomerOrder
	case model.PaymentTypeTopUp:
		reason = model.LedgerReasonTopUp
	}
	
	if err := s.sqlRepo.UpdateLedger(ctx, userID, reason, amount); err != nil {
		return fmt.Errorf("order.service.PaymentGatewayWebhook (UpdateLedger): %w", err)
	}
	
	return nil
}
