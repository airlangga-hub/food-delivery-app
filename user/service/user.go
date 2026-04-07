package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
)

type PaymentGatewayRepository interface {
	CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentGatewayResponse, error)
}

type userService struct {
	paymentGatewayRepository PaymentGatewayRepository
}

func NewUserService(userRepo PaymentGatewayRepository) *userService {
	return &userService{paymentGatewayRepository: userRepo}
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentGatewayResponse, error) {
	sessionLink, err := s.paymentGatewayRepository.CreatePaymentSession(ctx, userID, userEmail, amount)
	if err != nil {
		return model.PaymentGatewayResponse{}, fmt.Errorf("user.service.TopUpBalance: %w", err)
	}
	return sessionLink, nil
}
