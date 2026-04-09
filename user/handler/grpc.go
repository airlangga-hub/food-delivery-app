package handler

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/airlangga-hub/food-delivery-app/user/pb"
	"github.com/google/uuid"
)

type UserService interface {
	RegisterCustomer(ctx context.Context, input model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email string, password string) error
	GetUserInfo(ctx context.Context, email string) (model.UserInfo, error)
	PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error
	TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentLink, error)
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	Svc UserService
}

func NewHandler(svc UserService) *Handler {
	return &Handler{Svc: svc}
}

func (h *Handler) RegisterCustomer(context.Context, *pb.RegisterCustomerRequest) (*pb.RegisterCustomerResponse, error)

func (h *Handler) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)

func (h *Handler) GetUserInfo(context.Context, *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error)

func (h *Handler) TopUpBalance(context.Context, *pb.TopUpBalanceRequest) (*pb.TopUpBalanceResponse, error)

func (h *Handler) PaymentGatewayWebhook(context.Context, *pb.PaymentGatewayWebhookRequest) (*pb.PaymentGatewayWebhookResponse, error)

