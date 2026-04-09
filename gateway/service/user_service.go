package service

import (
	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	userpb "github.com/airlangga-hub/food-delivery-app/gateway/user_pb"
	"github.com/google/uuid"
)

type userService struct {
	userClient userpb.UserServiceClient
}

func NewUserService(userClient userpb.UserServiceClient) *userService {
	return &userService{userClient: userClient}
}

func (*userService) RegisterCustomer(user model.UserRegister) (model.UserInfo, error)
func (*userService) Login(email, password string) (string, error)
func (*userService) TopUpBalance(userID uuid.UUID, amount int) (model.PaymentLink, error)
func (*userService) GetUserInfo(userID uuid.UUID) (model.UserInfo, error)
