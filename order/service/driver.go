package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/order/model"
)

type DriverSQLRepository interface {
	DriverGetPendingOrders(ctx context.Context) ([]model.Order, error)
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