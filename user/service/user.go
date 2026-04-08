package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
)

type PaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type MongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type SQLRepository interface {
	UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error
}

type userService struct {
	paymentGatewayRepository PaymentGatewayRepository
	mongoRepository          MongoRepository
	sqlRepository            SQLRepository
}

func NewUserService(userRepo PaymentGatewayRepository, mongoRepo MongoRepository, sqlRepo SQLRepository) *userService {
	return &userService{paymentGatewayRepository: userRepo, mongoRepository: mongoRepo, sqlRepository: sqlRepo}
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentGatewayResponse, error) {
	items := []model.PaymentGatewayItem{
		{
			ReferenceID:   uuid.NewString(),
			Name:          "Top Up Customer Balance",
			Type:          "DIGITAL_SERVICE",
			Category:      "Top Up",
			NetUnitAmount: amount,
			Quantity:      1,
		},
	}

	paymentGatewayResp, err := s.paymentGatewayRepository.CreatePaymentSession(ctx, model.PaymentTypeTopUp, userID, userEmail, amount, items)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	if err := s.mongoRepository.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailStatus:            model.EmailStatusPending,
		PaymentType:            model.PaymentTypeTopUp,
		PaymentGatewayResponse: paymentGatewayResp,
	}); err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	return paymentGatewayResp, nil
}

func (s *userService) PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error {
	var reason model.LedgerReason
	switch paymentType {
	case model.PaymentTypeOrder:
		reason = model.LedgerReasonCustomerOrder
	case model.PaymentTypeTopUp:
		reason = model.LedgerReasonTopUp
	default:
		return fmt.Errorf("order.service.PaymentGatewayWebhook (invalid payment type: %s): %w", string(paymentType), model.ErrNotFound)
	}
	
	if err := s.sqlRepository.UpdateLedger(ctx, userID, reason, amount); err != nil {
		return fmt.Errorf("order.service.PaymentGatewayWebhook (UpdateLedger): %w", err)
	}
	
	return nil
}