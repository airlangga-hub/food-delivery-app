package model

import "go.mongodb.org/mongo-driver/v2/bson"

type (
	EmailStatus  string
	PaymentType  string
	LedgerReason string
)

const (
	EmailStatusPending EmailStatus = "PENDING"
	EmailStatusSending EmailStatus = "SENDING"
	EmailStatusSent    EmailStatus = "SENT"

	PaymentTypeTopUp PaymentType = "top_up"
	PaymentTypeOrder PaymentType = "order"

	LedgerReasonTopUp           LedgerReason = "customer_top_up"
	LedgerReasonCustomerOrder   LedgerReason = "customer_order"
	LedgerReasonDriverTakeOrder LedgerReason = "driver_take_order"
)

type PaymentRecord struct {
	ID                     bson.ObjectID          `bson:"_id,omitempty"`
	Email                  string                 `bson:"email"`
	EmailStatus            EmailStatus            `bson:"email_status"`
	PaymentType            PaymentType            `bson:"payment_type"`
	PaymentGatewayResponse PaymentGatewayResponse `bson:"payment_gateway_response"`
}
