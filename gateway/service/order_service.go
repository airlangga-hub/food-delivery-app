package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	orderpb "github.com/airlangga-hub/food-delivery-app/gateway/order_pb"
	pb "github.com/airlangga-hub/food-delivery-app/gateway/order_pb"
	"github.com/google/uuid"
)

type orderService struct {
	orderClient orderpb.OrderServiceClient
}

func NewOrderService(orderClient orderpb.OrderServiceClient) *orderService {
	return &orderService{orderClient: orderClient}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, deliveryFee int, items []model.Item) (model.Order, error) {
	itemsIn := make([]*pb.ItemsIn, len(items))

	for i, itm := range items {
		itemsIn[i] = &pb.ItemsIn{
			Id:       itm.ID,
			Quantity: int64(itm.Quantity),
		}
	}

	resp, err := s.orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:    userID.String(),
		UserEmail: userEmail,
		Order: &orderpb.OrderIn{
			DeliveryFee: int64(deliveryFee),
			ItemsIn:     itemsIn,
		},
	})

	if err != nil {
		return model.Order{}, fmt.Errorf("gateway.service.CreateOrder: %w", err)
	}

	restos := make([]model.Restaurant, len(resp.Restaurants))

	for i, resto := range resp.Restaurants {
		restoItems := make([]model.Item, len(resto.Items))

		for j, restoItem := range resto.Items {
			restoItems[j] = model.Item{
				ID:       restoItem.Id,
				Name:     restoItem.Name,
				Price:    int(restoItem.Price),
				Quantity: int(restoItem.Quantity),
			}
		}

		restos[i] = model.Restaurant{
			ID:      resto.Id,
			Name:    resto.Name,
			Address: resto.Address,
			Items:   restoItems,
		}
	}

	return model.Order{
		ID:                  resp.Id,
		Restaurants:         restos,
		DeliveryAddress:     resp.DeliveryAddress,
		CustomerPhoneNumber: resp.CustomerPhoneNumber,
		OrderStatus:         model.OrderStatus(resp.OrderStatus),
	}, nil
}

func (s *orderService) GetDrivers(orderID uuid.UUID) (model.FindDriver, error)

func (s *orderService) ChooseDriver(orderID, driverID uuid.UUID) (model.Order, error)

func (s *orderService) GetOrders(userID uuid.UUID) ([]model.Order, error)

func (s *orderService) GiveRating(orderID uuid.UUID) error

func (s *orderService) DriverGetPendingOrders() ([]model.Order, error)

func (s *orderService) DriverApplyToTakeOrder(driverID, orderID uuid.UUID) error

func (s *orderService) MarkOrderAsDone(orderID, driverID uuid.UUID) error
