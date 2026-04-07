package model

import "go.mongodb.org/mongo-driver/v2/bson"

type (
	EmailStatus string
)

const (
	EmailStatusPending EmailStatus = "PENDING"
	EmailStatusSending EmailStatus = "SENDING"
	EmailStatusSent    EmailStatus = "SENT"
)

type PaymentRecord struct {
	ID                     bson.ObjectID `bson:"_id,omitempty"`
	Email                  string        `bson:"email"`
	EmailSentStatus        EmailStatus   `bson:"email_sent_status"`
	PaymentType            string        `bson:"payment_type"`
	PaymentGatewayResponse any           `bson:"payment_gateway_response"`
}
