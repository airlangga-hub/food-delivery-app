package handler

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, input model.UserRegister) (model.UserInfo, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GetUserInfo(ctx context.Context, email string) (model.UserInfo, error)
	PaymentGatewayWebhook(ctx context.Context, userID uuid.UUID, paymentType model.PaymentType, amount int) error
	TopUpBalance(ctx context.Context, userID uuid.UUID, userEmail string, amount int) (model.PaymentLink, error)
}
