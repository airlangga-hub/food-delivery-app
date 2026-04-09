package repository

import (
	"context"

	"github.com/airlangga-hub/food-delivery-app/order/model"
	userpb "github.com/airlangga-hub/food-delivery-app/order/user_pb"
)

type mongoRepository struct {
	userClient userpb.UserServiceClient
}

func NewMongoRepository(userClient userpb.UserServiceClient) *mongoRepository {
	return &mongoRepository{userClient: userClient}
}

func (r *mongoRepository) CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error {
	r.userClient
}
