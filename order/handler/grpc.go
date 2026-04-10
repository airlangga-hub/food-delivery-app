package handler

import (
	"context"
	"errors"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/airlangga-hub/food-delivery-app/order/order_pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CustomerService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, order model.OrderIn) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
	ChooseDriver(ctx context.Context, orderID, driverID uuid.UUID) (model.Order, error)
	GiveRating(ctx context.Context, orderID uuid.UUID, rating int) error
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID, role model.RoleUser) ([]model.Order, error)
}

type DriverService interface {
	DriverGetPendingOrders(ctx context.Context) ([]model.Order, error)
	DriverApplyForOrder(ctx context.Context, orderID uuid.UUID, driverID uuid.UUID) error
	DriverCompleteOrder(ctx context.Context, orderID uuid.UUID, driverID uuid.UUID) error
}

type Handler struct {
	pb.UnimplementedOrderServiceServer
	customerSvc CustomerService
	driverSvc   DriverService
}

func NewHandler(customerSvc CustomerService, driverSvc DriverService) *Handler {
	return &Handler{customerSvc: customerSvc, driverSvc: driverSvc}
}

func (h *Handler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	items := make([]model.ItemsIn, len(req.Order.ItemsIn))

	for i, itm := range req.Order.ItemsIn {
		itmID, err := uuid.Parse(itm.Id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "order.handler.CreateOrder (Parse itm.Id): %v", err)
		}

		items[i] = model.ItemsIn{
			ID:       itmID,
			Quantity: int(itm.Quantity),
		}
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.CreateOrder (Parse UserId): %v", err)
	}

	order, err := h.customerSvc.CreateOrder(ctx, userID, req.UserEmail, model.OrderIn{
		DeliveryFee: int(req.Order.DeliveryFee),
		ItemsIn:     items,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.CreateOrder: %v", err)
	}

	restos := make([]*pb.Restaurant, len(order.Restaurants))

	for i, resto := range order.Restaurants {
		restoItems := make([]*pb.Item, len(resto.Items))

		for j, itm := range resto.Items {
			restoItems[j] = &pb.Item{
				Id:       itm.ID.String(),
				Name:     itm.Name,
				Price:    int64(itm.Price),
				Quantity: int64(itm.Quantity),
			}
		}

		restos[i] = &pb.Restaurant{
			Id:      resto.ID.String(),
			Name:    resto.Name,
			Address: resto.Address,
			Items:   restoItems,
		}
	}

	return &pb.CreateOrderResponse{
		Id:                  order.ID.String(),
		Restaurants:         restos,
		DeliveryAddress:     order.DeliveryAddress,
		CustomerPhoneNumber: order.CustomerPhoneNumber,
		OrderStatus:         string(order.OrderStatus),
		DeliveryFee:         int64(order.DeliveryFee),
		TotalFee:            int64(order.TotalFee),
		PaymentLink:         order.PaymentLink,
	}, nil
}

func (h *Handler) GetDrivers(ctx context.Context, req *pb.GetDriversRequest) (*pb.GetDriversResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.GetDrivers (Parse OrderId): %v", err)
	}

	drivers, err := h.customerSvc.GetDrivers(ctx, orderID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.GetDrivers (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.GetDrivers: %v", err)
	}

	resultDrivers := make([]*pb.Driver, len(drivers))

	for i, driver := range drivers {
		resultDrivers[i] = &pb.Driver{
			Id:            driver.ID.String(),
			AverageRating: driver.AverageRating,
			Name:          driver.Name,
			Bike:          driver.Bike,
			LicensePlate:  driver.LicensePlate,
			PhoneNumber:   driver.PhoneNumber,
		}
	}

	return &pb.GetDriversResponse{Drivers: resultDrivers}, nil
}

func (h *Handler) ChooseDriver(ctx context.Context, req *pb.ChooseDriverRequest) (*pb.ChooseDriverResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.ChooseDriver (Parse OrderId): %v", err)
	}

	driverID, err := uuid.Parse(req.DriverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.ChooseDriver (Parse DriverId): %v", err)
	}

	order, err := h.customerSvc.ChooseDriver(ctx, orderID, driverID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.ChooseDriver (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.ChooseDriver (Parse DriverId): %v", err)
	}

	restos := make([]*pb.Restaurant, len(order.Restaurants))

	for i, resto := range order.Restaurants {
		restoItems := make([]*pb.Item, len(resto.Items))

		for j, itm := range resto.Items {
			restoItems[j] = &pb.Item{
				Id:       itm.ID.String(),
				Name:     itm.Name,
				Price:    int64(itm.Price),
				Quantity: int64(itm.Quantity),
			}
		}

		restos[i] = &pb.Restaurant{
			Id:      resto.ID.String(),
			Name:    resto.Name,
			Address: resto.Address,
			Items:   restoItems,
		}
	}

	return &pb.ChooseDriverResponse{
		Id:                  order.ID.String(),
		Restaurants:         restos,
		DeliveryAddress:     order.DeliveryAddress,
		CustomerPhoneNumber: order.CustomerPhoneNumber,
		OrderStatus:         string(order.OrderStatus),
		Driver: &pb.Driver{
			Id:            order.Driver.ID.String(),
			AverageRating: order.Driver.AverageRating,
			Name:          order.Driver.Name,
			Bike:          order.Driver.Bike,
			LicensePlate:  order.Driver.LicensePlate,
			PhoneNumber:   order.Driver.PhoneNumber,
		},
		DeliveryFee: int64(order.DeliveryFee),
		TotalFee:    int64(order.TotalFee),
		PaymentLink: order.PaymentLink,
	}, nil
}

func (h *Handler) GiveRating(ctx context.Context, req *pb.GiveRatingRequest) (*pb.GiveRatingResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.GiveRating (Parse OrderId): %v", err)
	}

	if err := h.customerSvc.GiveRating(ctx, orderID, int(req.Rating)); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.GiveRating (orderID not found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.GiveRating: %v", err)
	}

	return nil, nil
}

func (h *Handler) DriverGetPendingOrders(ctx context.Context, req *emptypb.Empty) (*pb.DriverGetPendingOrdersResponse, error) {
	orders, err := h.driverSvc.DriverGetPendingOrders(ctx)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.DriverGetPendingOrders (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.DriverGetPendingOrders: %v", err)
	}

	resultOrders := make([]*pb.Order, len(orders))

	for i, order := range orders {
		restos := make([]*pb.Restaurant, len(order.Restaurants))

		for j, resto := range order.Restaurants {
			restoItems := make([]*pb.Item, len(resto.Items))

			for k, restoItem := range resto.Items {
				restoItems[k] = &pb.Item{
					Id:       restoItem.ID.String(),
					Name:     restoItem.Name,
					Price:    int64(restoItem.Price),
					Quantity: int64(restoItem.Quantity),
				}
			}

			restos[j] = &pb.Restaurant{
				Id:      resto.ID.String(),
				Name:    resto.Name,
				Address: resto.Address,
				Items:   restoItems,
			}
		}

		resultOrders[i] = &pb.Order{
			Id:                  order.ID.String(),
			Restaurants:         restos,
			DeliveryAddress:     order.DeliveryAddress,
			CustomerPhoneNumber: order.CustomerPhoneNumber,
			OrderStatus:         string(order.OrderStatus),
			DeliveryFee:         int64(order.DeliveryFee),
			TotalFee:            int64(order.TotalFee),
		}
	}

	return &pb.DriverGetPendingOrdersResponse{Orders: resultOrders}, nil
}

func (h *Handler) DriverApplyForOrder(ctx context.Context, req *pb.DriverApplyForOrderRequest) (*pb.DriverApplyForOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.DriverApplyForOrder (Parse OrderId): %v", err)
	}

	driverID, err := uuid.Parse(req.DriverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.DriverApplyForOrder (Parse DriverId): %v", err)
	}

	if err := h.driverSvc.DriverApplyForOrder(ctx, orderID, driverID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.DriverApplyForOrder (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.DriverApplyForOrder: %v", err)
	}

	return nil, nil
}

func (h *Handler) DriverCompleteOrder(ctx context.Context, req *pb.DriverCompleteOrderRequest) (*pb.DriverCompleteOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.DriverCompleteOrder (Parse OrderId): %v", err)
	}

	driverID, err := uuid.Parse(req.DriverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.DriverCompleteOrder (Parse DriverId): %v", err)
	}

	if err := h.driverSvc.DriverCompleteOrder(ctx, orderID, driverID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.DriverCompleteOrder (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.DriverCompleteOrder: %v", err)
	}

	return nil, nil
}

func (h *Handler) CreatePaymentSession(ctx context.Context, req *pb.CreatePaymentSessionRequest) (*pb.CreatePaymentSessionResponse, error) {
	items := make([]model.PaymentGatewayItem, len(req.Items))

	for i, itm := range req.Items {
		items[i] = model.PaymentGatewayItem{
			ReferenceID:   itm.ReferenceId,
			Name:          itm.Name,
			Description:   itm.Description,
			Type:          itm.Type,
			Category:      itm.Category,
			NetUnitAmount: int(itm.NetUnitAmount),
			Quantity:      int(itm.Quantity),
			URL:           itm.Url,
		}
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.CreatePaymentSession (Parse UserId): %v", err)
	}

	pgResp, err := h.customerSvc.CreatePaymentSession(ctx, model.PaymentType(req.PaymentType), userID, req.UserEmail, int(req.Amount), items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.CreatePaymentSession: %v", err)
	}

	pgItems := make([]*pb.PaymentGatewayItem, len(pgResp.Items))

	for i, pgItm := range pgResp.Items {
		pgItems[i] = &pb.PaymentGatewayItem{
			ReferenceId:   pgItm.ReferenceID,
			Name:          pgItm.Name,
			Description:   pgItm.Description,
			Type:          pgItm.Type,
			Category:      pgItm.Category,
			NetUnitAmount: int64(pgItm.NetUnitAmount),
			Quantity:      int64(pgItm.Quantity),
			Url:           pgItm.URL,
		}
	}

	return &pb.CreatePaymentSessionResponse{
		PaymentSessionId: pgResp.PaymentSessionID,
		Created:          pgResp.Created,
		Updated:          pgResp.Updated,
		Status:           pgResp.Status,
		ReferenceId:      pgResp.ReferenceID,
		Currency:         pgResp.Currency,
		Amount:           pgResp.Amount,
		Country:          pgResp.Country,
		ExpiresAt:        pgResp.ExpiresAt,
		SessionType:      pgResp.SessionType,
		Mode:             pgResp.Mode,
		Locale:           pgResp.Locale,
		BusinessId:       pgResp.BusinessID,
		CustomerId:       pgResp.CustomerID,
		CaptureMethod:    pgResp.CaptureMethod,
		Description:      pgResp.Description,
		Items:            pgItems,
		SuccessReturnUrl: pgResp.SuccessReturnURL,
		CancelReturnUrl:  pgResp.CancelReturnURL,
		PaymentLinkUrl:   pgResp.PaymentLinkURL,
	}, nil
}

func (h *Handler) GetOrdersByUserID(ctx context.Context, req *pb.GetOrdersByUserIDRequest) (*pb.GetOrdersByUserIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "order.handler.GetOrdersByUserID (Parse UserId): %v", err)
	}

	orders, err := h.customerSvc.GetOrdersByUserID(ctx, userID, model.RoleUser(req.Role))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "order.handler.GetOrdersByUserID (no rows found): %v", err)
		}
		return nil, status.Errorf(codes.Internal, "order.handler.GetOrdersByUserID: %v", err)
	}

	resultOrders := make([]*pb.Order, len(orders))

	for i, order := range orders {
		restos := make([]*pb.Restaurant, len(order.Restaurants))

		for j, resto := range order.Restaurants {
			restoItems := make([]*pb.Item, len(resto.Items))

			for k, itm := range resto.Items {
				restoItems[k] = &pb.Item{
					Id:       itm.ID.String(),
					Name:     itm.Name,
					Price:    int64(itm.Price),
					Quantity: int64(itm.Quantity),
				}
			}

			restos[j] = &pb.Restaurant{
				Id:      resto.ID.String(),
				Name:    resto.Name,
				Address: resto.Address,
				Items:   restoItems,
			}
		}

		resultOrders[i] = &pb.Order{
			Id:                  order.ID.String(),
			Restaurants:         restos,
			DeliveryAddress:     order.DeliveryAddress,
			CustomerPhoneNumber: order.CustomerPhoneNumber,
			OrderStatus:         string(order.OrderStatus),
			DeliveryFee:         int64(order.DeliveryFee),
			TotalFee:            int64(order.TotalFee),
			PaymentLink:         order.PaymentLink,
			Driver: &pb.Driver{
				AverageRating: order.Driver.AverageRating,
				Name:          order.Driver.Name,
				Bike:          order.Driver.Bike,
				LicensePlate:  order.Driver.LicensePlate,
				PhoneNumber:   order.Driver.PhoneNumber,
			},
		}
	}
	
	return &pb.GetOrdersByUserIDResponse{Orders: resultOrders}, nil
}
