package service

import (
	"context"
)

type UserRepository interface {
	CreatePaymentSession(ctx context.Context, amount int) (string, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *userService {
	return &userService{userRepo: userRepo}
}
