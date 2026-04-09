package service

import (
	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	orderpb "github.com/airlangga-hub/food-delivery-app/gateway/order_pb"
	"github.com/google/uuid"
)

type orderService struct {
	orderClient orderpb.OrderServiceClient
}

func NewOrderService(orderClient orderpb.OrderServiceClient) *orderService {
	return &orderService{orderClient: orderClient}
}

func (s *orderService) CreateOrder(userID uuid.UUID, items []model.Item) (model.Order, error)
func (s *orderService) GetDrivers(orderID uuid.UUID) (model.FindDriver, error)
func (s *orderService) ChooseDriver(orderID, driverID uuid.UUID) (model.Order, error)
func (s *orderService) GetOrders(userID uuid.UUID) ([]model.Order, error)
func (s *orderService) GiveRating(orderID uuid.UUID) error
func (s *orderService) DriverGetPendingOrders() ([]model.Order, error)
func (s *orderService) DriverApplyToTakeOrder(driverID, orderID uuid.UUID) error
func (s *orderService) MarkOrderAsDone(orderID, driverID uuid.UUID) error
