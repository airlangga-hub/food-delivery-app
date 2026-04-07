package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreatePaymentSession(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (string, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (string, error) {
	sessionLink, err := s.userRepo.CreatePaymentSession(ctx, userID, userEmail, amount)
	if err != nil {
		return "", fmt.Errorf("user.service.TopUpBalance: %w", err)
	}
	return sessionLink, nil
}
