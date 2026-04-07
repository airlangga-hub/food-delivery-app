package service

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, order model.Order) (model.Order, error)
}

type userService struct {
	customerRepo CustomerRepository
}

func NewUserService(customerRepo CustomerRepository) *userService {
	return &userService{customerRepo: customerRepo}
}
