package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	orderpb "github.com/airlangga-hub/food-delivery-app/gateway/order_pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderService struct {
	orderClient orderpb.OrderServiceClient
	logger      *slog.Logger
}

func NewOrderService(orderClient orderpb.OrderServiceClient, logger *slog.Logger) *orderService {
	return &orderService{orderClient: orderClient, logger: logger}
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
			slog.Info("gateway.service.CreateOrder (FromError): not gRPC error: %w", slog.Any("error", err))
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
		DeliveryFee:         int(resp.DeliveryFee),
		TotalFee:            int(resp.TotalFee),
		PaymentLink:         resp.PaymentLink,
	}, nil
}

func (s *orderService) GetDrivers(ctx context.Context, orderID string) (model.FindDriver, error) {
	resp, err := s.orderClient.GetDrivers(ctx, &orderpb.GetDriversRequest{OrderId: orderID})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.GetDrivers (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return model.FindDriver{}, fmt.Errorf("gateway.service.GetDrivers (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return model.FindDriver{}, fmt.Errorf("gateway.service.GetDrivers: %w", err)
	}

	drivers := make([]model.Driver, len(resp.Drivers))

	for i, drv := range resp.Drivers {
		drivers[i] = model.Driver{
			ID:            drv.Id,
			AverageRating: drv.AverageRating,
			Name:          drv.Name,
			Bike:          drv.Bike,
			LicensePlate:  drv.LicensePlate,
			PhoneNumber:   drv.PhoneNumber,
		}
	}

	return model.FindDriver{OrderApplicants: drivers}, nil
}

func (s *orderService) ChooseDriver(ctx context.Context, orderID, driverID string) (model.Order, error) {
	resp, err := s.orderClient.ChooseDriver(ctx, &orderpb.ChooseDriverRequest{OrderId: orderID, DriverId: driverID})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.GetDrivers (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return model.Order{}, fmt.Errorf("gateway.service.GetDrivers (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return model.Order{}, fmt.Errorf("gateway.service.GetDrivers: %w", err)
	}

	restos := make([]model.Restaurant, len(resp.Restaurants))

	for i, resto := range resp.Restaurants {
		restoItems := make([]model.Item, len(resto.Items))

		for j, itm := range resto.Items {
			restoItems[j] = model.Item{
				ID:       itm.Id,
				Name:     itm.Name,
				Price:    int(itm.Price),
				Quantity: int(itm.Quantity),
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
		DeliveryFee:         int(resp.DeliveryFee),
		TotalFee:            int(resp.TotalFee),
		PaymentLink:         resp.PaymentLink,
		Driver: &model.Driver{
			ID:            resp.Driver.Id,
			AverageRating: resp.Driver.AverageRating,
			Name:          resp.Driver.Name,
			Bike:          resp.Driver.Bike,
			LicensePlate:  resp.Driver.LicensePlate,
			PhoneNumber:   resp.Driver.PhoneNumber,
		},
	}, nil
}

func (s *orderService) CustomerGetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	resp, err := s.orderClient.GetOrdersByUserID(ctx, &orderpb.GetOrdersByUserIDRequest{UserId: userID, Role: string(model.RoleUserCustomer)})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.GetOrders (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return nil, fmt.Errorf("gateway.service.GetOrders (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return nil, fmt.Errorf("gateway.service.GetOrders: %w", err)
	}

	resultOrders := make([]model.Order, len(resp.Orders))

	for i, order := range resp.Orders {
		restos := make([]model.Restaurant, len(order.Restaurants))

		for j, resto := range order.Restaurants {
			restoItems := make([]model.Item, len(resto.Items))

			for k, itm := range resto.Items {
				restoItems[k] = model.Item{
					ID:       itm.Id,
					Name:     itm.Name,
					Price:    int(itm.Price),
					Quantity: int(itm.Quantity),
				}
			}

			restos[j] = model.Restaurant{
				ID:      resto.Id,
				Name:    resto.Name,
				Address: resto.Address,
				Items:   restoItems,
			}
		}

		resultOrders[i] = model.Order{
			ID:                  order.Id,
			Restaurants:         restos,
			DeliveryAddress:     order.DeliveryAddress,
			CustomerPhoneNumber: order.CustomerPhoneNumber,
			OrderStatus:         model.OrderStatus(order.OrderStatus),
			DeliveryFee:         int(order.DeliveryFee),
			TotalFee:            int(order.TotalFee),
			PaymentLink:         order.PaymentLink,
			Driver: &model.Driver{
				AverageRating: order.Driver.AverageRating,
				Name:          order.Driver.Name,
				Bike:          order.Driver.Bike,
				LicensePlate:  order.Driver.LicensePlate,
				PhoneNumber:   order.Driver.PhoneNumber,
			},
		}
	}
	
	return resultOrders, nil
}

func (s *orderService) GiveRating(ctx context.Context, orderID string, rating int) error {
	_, err := s.orderClient.GiveRating(ctx, &orderpb.GiveRatingRequest{OrderId: orderID, Rating: int64(rating)})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.GetOrders (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return fmt.Errorf("gateway.service.GetOrders (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return fmt.Errorf("gateway.service.GetOrders: %w", err)
	}
	
	return nil
}

func (s *orderService) DriverGetPendingOrders(ctx context.Context) ([]model.Order, error) {
	resp, err := s.orderClient.DriverGetPendingOrders(ctx, &emptypb.Empty{})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.DriverGetPendingOrders (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return nil, fmt.Errorf("gateway.service.DriverGetPendingOrders (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return nil, fmt.Errorf("gateway.service.DriverGetPendingOrders: %w", err)
	}
	
	resultOrders := make([]model.Order, len(resp.Orders))

	for i, order := range resp.Orders {
		restos := make([]model.Restaurant, len(order.Restaurants))

		for j, resto := range order.Restaurants {
			restoItems := make([]model.Item, len(resto.Items))

			for k, restoItem := range resto.Items {
				restoItems[k] = model.Item{
					ID:       restoItem.Id,
					Name:     restoItem.Name,
					Price:    int(restoItem.Price),
					Quantity: int(restoItem.Quantity),
				}
			}

			restos[j] = model.Restaurant{
				ID:      resto.Id,
				Name:    resto.Name,
				Address: resto.Address,
				Items:   restoItems,
			}
		}

		resultOrders[i] = model.Order{
			ID:                  order.Id,
			Restaurants:         restos,
			DeliveryAddress:     order.DeliveryAddress,
			CustomerPhoneNumber: order.CustomerPhoneNumber,
			OrderStatus:         model.OrderStatus(order.OrderStatus),
			DeliveryFee:         int(order.DeliveryFee),
			TotalFee:            int(order.TotalFee),
		}
	}
	
	return resultOrders, nil
}

func (s *orderService) DriverApplyToTakeOrder(ctx context.Context, driverID, orderID string) error {
	_, err := s.orderClient.DriverApplyForOrder(ctx, &orderpb.DriverApplyForOrderRequest{OrderId: orderID, DriverId: driverID})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.DriverApplyToTakeOrder (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return fmt.Errorf("gateway.service.DriverApplyToTakeOrder (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return fmt.Errorf("gateway.service.DriverApplyToTakeOrder: %w", err)
	}
	
	return nil
}

func (s *orderService) MarkOrderAsDone(ctx context.Context, orderID, driverID string) error
