package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	userpb "github.com/airlangga-hub/food-delivery-app/gateway/user_pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	userClient userpb.UserServiceClient
}

func NewUserService(userClient userpb.UserServiceClient) *userService {
	return &userService{userClient: userClient}
}

func (s *userService) RegisterCustomer(ctx context.Context, user model.UserRegister) (model.UserInfo, error) {
	resp, err := s.userClient.RegisterCustomer(ctx, &userpb.RegisterCustomerRequest{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.RegisterCustomer (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return model.UserInfo{}, fmt.Errorf("gateway.service.RegisterCustomer (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return model.UserInfo{}, fmt.Errorf("gateway.service.RegisterCustomer: %w", err)
	}

	return model.UserInfo{
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		Address:   resp.Address,
		Balance:   int(resp.Balance),
	}, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := s.userClient.Login(ctx, &userpb.LoginRequest{Email: email, Password: password})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.Login (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return "", fmt.Errorf("gateway.service.Login (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return "", fmt.Errorf("gateway.service.Login: %w", err)
	}
	
	return resp.Token, nil
}

func (s *userService) TopUpBalance(ctx context.Context, userID, userEmail string, amount int) (model.PaymentLink, error) {
	resp, err := s.userClient.TopUpBalance(ctx, &userpb.TopUpBalanceRequest{UserId: userID, UserEmail: userEmail, Amount: int64(amount)})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.TopUpBalance (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return model.PaymentLink{}, fmt.Errorf("gateway.service.TopUpBalance (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return model.PaymentLink{}, fmt.Errorf("gateway.service.TopUpBalance: %w", err)
	}
	
	return model.PaymentLink{PaymentLink: resp.PaymentLink}, nil
}

func (s *userService) GetUserInfo(ctx context.Context, userID string) (model.UserInfo, error)

func (s *userService) PaymentGatewayWebhook(ctx context.Context, userID string, paymentType model.PaymentType, amount int) error {
	_, err := s.userClient.PaymentGatewayWebhook(ctx, &userpb.PaymentGatewayWebhookRequest{UserId: userID, PaymentType: string(paymentType), Amount: int64(amount)})
	
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			slog.Info("gateway.service.PaymentGatewayWebhook (FromError): not gRPC error: %w", slog.Any("error", err))
		}

		if st.Code() == codes.NotFound {
			return fmt.Errorf("gateway.service.PaymentGatewayWebhook (no rows found): %w: %w", model.ErrNotFound, err)
		}

		return fmt.Errorf("gateway.service.PaymentGatewayWebhook: %w", err)
	}
	
	return nil
}