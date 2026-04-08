package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type CustomerSQLRepository interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, order model.OrderIn) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
}

type CustomerPaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type CustomerMongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type customerService struct {
	customerSqlRepo            CustomerSQLRepository
	customerPaymentGatewayRepo CustomerPaymentGatewayRepository
	customerMongoRepo          CustomerMongoRepository
}

func NewCustomerService(customerRepo CustomerSQLRepository, paymentGatewayRepo CustomerPaymentGatewayRepository, mongoRepo CustomerMongoRepository) *customerService {
	return &customerService{customerSqlRepo: customerRepo, customerPaymentGatewayRepo: paymentGatewayRepo, customerMongoRepo: mongoRepo}
}

func (s *customerService) CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, order model.OrderIn) (model.Order, error) {
	oorder, err := s.customerSqlRepo.CreateOrder(ctx, userID, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder (sqlRepo.CreateOrder): %w", err)
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

	paymentGatewayResponse, err := s.customerPaymentGatewayRepo.CreatePaymentSession(ctx, model.PaymentTypeOrder, userID, userEmail, oorder.TotalFee, items)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder (CreatePaymentSession): %w", err)
	}

	if err := s.customerMongoRepo.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailStatus:            model.EmailStatusPending,
		PaymentType:            model.PaymentTypeOrder,
		PaymentGatewayResponse: paymentGatewayResponse,
	}); err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder (CreatePaymentRecord): %w", err)
	}

	oorder.PaymentLink = paymentGatewayResponse.PaymentLinkURL

	return oorder, nil
}

func (s *customerService) GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error) {
	drivers, err := s.customerSqlRepo.GetDrivers(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order.customer_service.GetDrivers: %w", err)
	}
	return drivers, nil
}
