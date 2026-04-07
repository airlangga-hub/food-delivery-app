package repository

import "go.mongodb.org/mongo-driver/v2/bson"

type PaymentRecord struct {
	ID                     bson.ObjectID `bson:"_id,omitempty"`
	Email                  string        `bson:"email"`
	EmailSentStatus        EmailStatus   `bson:"email_sent_status"`
	PaymentGatewayResponse any           `bson:"payment_gateway_response"`
}
