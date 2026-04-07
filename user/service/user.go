package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type XenditRepository interface {
	CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (string, error)
}

type userService struct {
	xenditRepo XenditRepository
}

func NewUserService(userRepo XenditRepository) *userService {
	return &userService{xenditRepo: userRepo}
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (string, error) {
	sessionLink, err := s.xenditRepo.CreatePaymentSession(ctx, userID, userEmail, amount)
	if err != nil {
		return "", fmt.Errorf("user.service.TopUpBalance: %w", err)
	}
	return sessionLink, nil
}
