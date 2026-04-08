package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/google/uuid"
)

type DriverSQLRepository interface {
	DriverGetPendingOrders(ctx context.Context) ([]model.Order, error)
	DriverApplyForOrder(ctx context.Context, orderID uuid.UUID, driverID uuid.UUID) error
}

type driverService struct {
	driverSqlRepo DriverSQLRepository
}

func NewDriverSQLRepository(driverSqlRepo DriverSQLRepository) *driverService {
	return &driverService{driverSqlRepo: driverSqlRepo}
}

func (s *driverService) DriverGetPendingOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := s.driverSqlRepo.DriverGetPendingOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("order.service.DriverGetPendingOrders: %w", err)
	}
	return orders, nil
}

func (s *driverService) DriverApplyForOrder(ctx context.Context, orderID uuid.UUID, driverID uuid.UUID) error {
	if err := s.driverSqlRepo.DriverApplyForOrder(ctx, orderID, driverID); err != nil {
		return fmt.Errorf("order.service.DriverApplyForOrder: %w", err)
	}
	return nil
}