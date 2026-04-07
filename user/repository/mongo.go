package repository

import (
	"context"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoRepository struct {
	paymentsCol *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *mongoRepository {
	paymentsCol := db.Collection("payments")
	return &mongoRepository{paymentsCol: paymentsCol}
}

func (r *mongoRepository) CreatePaymentRecord(ctx context.Context, paymentRecord model.PaymentRecord) error {
	_, err := r.paymentsCol.InsertOne(ctx, paymentRecord)
	if err != nil {
		return fmt.Errorf("user.repository.CreatePaymentRecord: %w", err)
	}
	return nil
}
