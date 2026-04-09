package handler

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	"github.com/airlangga-hub/food-delivery-app/order/pb"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type CustomerService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, userEmail string, order model.OrderIn) (model.Order, error)
	GetDrivers(ctx context.Context, orderID uuid.UUID) ([]model.Driver, error)
	ChooseDriver(ctx context.Context, orderID, driverID uuid.UUID) (model.Order, error)
	GiveRating(ctx context.Context, orderID uuid.UUID, rating int) error
	CreatePaymentSession(ctx context.Context, paymentType model.PaymentType, userID uuid.UUID, userEmail string, amount int, items []model.PaymentGatewayItem) (model.PaymentGatewayResponse, error)
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
	
}

func (h *Handler) GetDrivers(ctx context.Context, req *pb.GetDriversRequest) (*pb.GetDriversResponse, error)

func (h *Handler) ChooseDriver(ctx context.Context, req *pb.ChooseDriverRequest) (*pb.ChooseDriverResponse, error)

func (h *Handler) GiveRating(ctx context.Context, req *pb.GiveRatingRequest) (*pb.GiveRatingResponse, error)

func (h *Handler) DriverGetPendingOrders(ctx context.Context, req *emptypb.Empty) (*pb.DriverGetPendingOrdersResponse, error)

func (h *Handler) DriverApplyForOrder(ctx context.Context, req *pb.DriverApplyForOrderRequest) (*pb.DriverApplyForOrderResponse, error)

func (h *Handler) DriverCompleteOrder(ctx context.Context, req *pb.DriverCompleteOrderRequest) (*pb.DriverCompleteOrderResponse, error)

func (h *Handler) CreatePaymentSession(ctx context.Context, req *pb.CreatePaymentSessionRequest) (*pb.CreatePaymentSessionResponse, error)
