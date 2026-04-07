package model

import "go.mongodb.org/mongo-driver/v2/bson"

type (
	EmailStatus string
	PaymentType string
)

const (
	EmailStatusPending EmailStatus = "PENDING"
	EmailStatusSending EmailStatus = "SENDING"
	EmailStatusSent    EmailStatus = "SENT"

	PaymentTypeTopUp PaymentType = "top_up"
	PaymentTypeOrder PaymentType = "order"
)

type PaymentRecord struct {
	ID                     bson.ObjectID `bson:"_id,omitempty"`
	Email                  string        `bson:"email"`
	EmailStatus            EmailStatus   `bson:"email_status"`
	PaymentType            PaymentType   `bson:"payment_type"`
	PaymentGatewayResponse any           `bson:"payment_gateway_response"`
}
