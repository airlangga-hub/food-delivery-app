package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type CustomerSQLRepository interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, order model.OrderIn) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
	ChooseDriver(ctx context.Context, orderID, driverID uuid.UUID) (model.Order, error)
	GiveRating(ctx context.Context, orderID uuid.UUID, rating int) error
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID, role model.RoleUser) ([]model.Order, error)
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

func NewCustomerService(customerSqlRepo CustomerSQLRepository, customerPaymentGatewayRepo CustomerPaymentGatewayRepository, customerMongoRepo CustomerMongoRepository) *customerService {
	return &customerService{customerSqlRepo: customerSqlRepo, customerPaymentGatewayRepo: customerPaymentGatewayRepo, customerMongoRepo: customerMongoRepo}
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

	items = append(items, model.PaymentGatewayItem{
		ReferenceID:   uuid.NewString(),
		Name:          "Delivery Fee",
		Type:          "PHYSICAL_SERVICE",
		Category:      "Order",
		NetUnitAmount: order.DeliveryFee,
		Quantity:      1,
	})

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

func (s *customerService) ChooseDriver(ctx context.Context, orderID, driverID uuid.UUID) (model.Order, error) {
	order, err := s.customerSqlRepo.ChooseDriver(ctx, orderID, driverID)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.service.ChooseDriver (customerSqlRepo.ChooseDriver): %w", err)
	}
	return order, nil
}

func (s *customerService) GiveRating(ctx context.Context, orderID uuid.UUID, rating int) error {
	if err := s.customerSqlRepo.GiveRating(ctx, orderID, rating); err != nil {
		return fmt.Errorf("order.service.GiveRating: %w", err)
	}
	return nil
}

func (s *customerService) CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error) {
	paymentGatewayResp, err := s.customerPaymentGatewayRepo.CreatePaymentSession(ctx, paymentType, userID, userEmail, amount, items)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("order.service.CreatePaymentSession: %w", err)
	}
	return paymentGatewayResp, nil
}

func (s *customerService) GetOrdersByUserID(ctx context.Context, userID uuid.UUID, role model.RoleUser) ([]model.Order, error) {
	orders, err := s.customerSqlRepo.GetOrdersByUserID(ctx, userID, role)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("order.service.GetOrdersByUserID (no rows found): %w", err)
		}
		return nil, fmt.Errorf("order.service.GetOrdersByUserID: %w", err)
	}
	return orders, nil
}
