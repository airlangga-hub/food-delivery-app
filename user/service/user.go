package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
)

type PaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
}

type MongoRepository interface {
	CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error
}

type userService struct {
	paymentGatewayRepository PaymentGatewayRepository
	mongoRepository          MongoRepository
}

func NewUserService(userRepo PaymentGatewayRepository) *userService {
	return &userService{paymentGatewayRepository: userRepo}
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

	paymentGatewayResp, err := s.paymentGatewayRepository.CreatePaymentSession(ctx, userID, userEmail, amount, items)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	if err := s.mongoRepository.CreatePaymentRecord(ctx, model.PaymentRecord{
		Email:                  userEmail,
		EmailSentStatus:        model.EmailStatusPending,
		PaymentType:            model.PaymentTypeTopUp,
		PaymentGatewayResponse: paymentGatewayResp,
	}); err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}

	return paymentGatewayResp, nil
}
