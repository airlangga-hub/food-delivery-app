package service

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	orderpb "github.com/airlangga-hub/food-delivery-app/gateway/order_pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	orderClient orderpb.OrderServiceClient
}

func NewOrderService(orderClient orderpb.OrderServiceClient) *orderService {
	return &orderService{orderClient: orderClient}
}

func (s *orderService) CreateOrder(ctx context.Context, userID string, userEmail string, deliveryFee int, items []model.Item) (model.Order, error) {
	itemsIn := make([]*orderpb.ItemsIn, len(items))

	for i, itm := range items {
		itemsIn[i] = &orderpb.ItemsIn{
			Id:       itm.ID,
			Quantity: int64(itm.Quantity),
		}
	}

	resp, err := s.orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:    userID,
		UserEmail: userEmail,
		Order: &orderpb.OrderIn{
			DeliveryFee: int64(deliveryFee),
			ItemsIn:     itemsIn,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return model.Order{}, fmt.Errorf("gateway.service.CreateOrder (FromError): not gRPC error: %w", err)
		}

		if st.Code() == codes.NotFound {
			return model.Order{}, fmt.Errorf("gateway.service.CreateOrder (no rows found): %w: %w", model.ErrNotFound, err)
		}

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

func (s *orderService) GetDrivers(ctx context.Context, orderID string) (model.FindDriver, error)

func (s *orderService) ChooseDriver(ctx context.Context, orderID, driverID string) (model.Order, error)

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]model.Order, error)

func (s *orderService) GiveRating(ctx context.Context, orderID string) error

func (s *orderService) DriverGetPendingOrders(ctx context.Context) ([]model.Order, error)

func (s *orderService) DriverApplyToTakeOrder(ctx context.Context, driverID, orderID string) error

func (s *orderService) MarkOrderAsDone(ctx context.Context, orderID, driverID string) error
