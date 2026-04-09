package model

import "go.mongodb.org/mongo-driver/v2/bson"

type PaymentRecord struct {
	ID                     bson.ObjectID          `bson:"_id,omitempty"`
	Email                  string                 `bson:"email"`
	EmailStatus            EmailStatus            `bson:"email_status"`
	PaymentType            PaymentType            `bson:"payment_type"`
	PaymentGatewayResponse PaymentGatewayResponse `bson:"payment_gateway_response"`
}
