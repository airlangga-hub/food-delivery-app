package service

import (
	"context"
	"fmt"

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

func (s *userService) CreateOrder(ctx context.Context, userID uuid.UUID, order model.Order) (model.Order, error) {
	order, err := s.customerRepo.CreateOrder(ctx, userID, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("order.customer_service.CreateOrder: %w", err)
	}
	return order, nil
}
