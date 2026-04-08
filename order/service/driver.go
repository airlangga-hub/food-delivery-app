package service

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/order/model"
)

type DriverSQLRepository interface {
	GetPendingOrders(ctx context.Context) ([]model.Order, error)
}

type driverService struct {
	driverSqlRepo DriverSQLRepository
}

func NewDriverSQLRepository(driverSqlRepo DriverSQLRepository) *driverService {
	return &driverService{driverSqlRepo: driverSqlRepo}
}