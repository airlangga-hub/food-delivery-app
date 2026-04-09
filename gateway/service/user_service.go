package service

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	userpb "github.com/airlangga-hub/food-delivery-app/gateway/user_pb"
)

type userService struct {
	userClient userpb.UserServiceClient
}

func NewUserService(userClient userpb.UserServiceClient) *userService {
	return &userService{userClient: userClient}
}

func (*userService) RegisterCustomer(ctx context.Context, user model.UserRegister) (model.UserInfo, error)
func (*userService) Login(ctx context.Context, email, password string) (string, error)
func (*userService) TopUpBalance(ctx context.Context, userID string, amount int) (model.PaymentLink, error)
func (*userService) GetUserInfo(ctx context.Context, userID string) (model.UserInfo, error)
