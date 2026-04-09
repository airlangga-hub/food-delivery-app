package handler

import (
	"context"
	"errors"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/airlangga-hub/food-delivery-app/user/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *Handler) RegisterCustomer(ctx context.Context, req *pb.RegisterCustomerRequest) (*pb.RegisterCustomerResponse, error) {
	userinfo, err := h.Svc.RegisterCustomer(ctx, model.UserRegister{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "user.handler.RegisterCustomer: %v", err)
	}

	return &pb.RegisterCustomerResponse{
		FirstName:   userinfo.FirstName,
		LastName:    userinfo.LastName,
		Email:       userinfo.Email,
		Address:     userinfo.Address,
		PhoneNumber: userinfo.PhoneNumber,
		Balance:     int64(userinfo.Balance),
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := h.Svc.Login(ctx, req.Email, req.Password); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "user.handler.Login (user not found): %v: %v", model.ErrNotFound, err)
		}
		return nil, status.Errorf(codes.Internal, "user.handler.Login: %v", err)
	}
	return nil, nil
}

func (h *Handler) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userinfo, err := h.Svc.GetUserInfo(ctx, req.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.Internal, "user.handler.GetUserInfo (user not found): %v: %v", model.ErrNotFound, err)
		}
		return nil, status.Errorf(codes.Internal, "user.handler.GetUserInfo: %v", err)
	}
	return &pb.GetUserInfoResponse{
		FirstName:   userinfo.FirstName,
		LastName:    userinfo.LastName,
		Email:       userinfo.Email,
		Address:     userinfo.Address,
		PhoneNumber: userinfo.PhoneNumber,
		Balance:     int64(userinfo.Balance),
	}, nil
}

func (h *Handler) PaymentGatewayWebhook(ctx context.Context, req *pb.PaymentGatewayWebhookRequest) (*pb.PaymentGatewayWebhookResponse, error)

func (h *Handler) TopUpBalance(ctx context.Context, req *pb.TopUpBalanceRequest) (*pb.TopUpBalanceResponse, error)
