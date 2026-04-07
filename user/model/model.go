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
	ID                     bson.ObjectID          `bson:"_id,omitempty"`
	Email                  string                 `bson:"email"`
	EmailSentStatus        EmailStatus            `bson:"email_sent_status"`
	PaymentType            PaymentType            `bson:"payment_type"`
	PaymentGatewayResponse PaymentGatewayResponse `bson:"payment_gateway_response"`
}

type PaymentGatewayResponse struct {
	PaymentSessionID string               `bson:"payment_session_id"`
	Created          string               `bson:"created"`
	Updated          string               `bson:"updated"`
	Status           string               `bson:"status"`
	ReferenceID      string               `bson:"reference_id"`
	Currency         string               `bson:"currency"`
	Amount           int64                `bson:"amount"`
	Country          string               `bson:"country"`
	ExpiresAt        string               `bson:"expires_at"`
	SessionType      string               `bson:"session_type"`
	Mode             string               `bson:"mode"`
	Locale           string               `bson:"locale"`
	BusinessID       string               `bson:"business_id"`
	CustomerID       string               `bson:"customer_id"`
	CaptureMethod    string               `bson:"capture_method"`
	Description      string               `bson:"description"`
	Items            []PaymentGatewayItem `bson:"items"`
	SuccessReturnURL string               `bson:"success_return_url"`
	CancelReturnURL  string               `bson:"cancel_return_url"`
	PaymentLinkURL   string               `bson:"payment_link_url"`
}

type PaymentGatewayItem struct {
	ReferenceID   string `bson:"reference_id"`
	Name          string `bson:"name"`
	Description   string `bson:"description"`
	Type          string `bson:"type"`
	Category      string `bson:"category"`
	NetUnitAmount int    `bson:"net_unit_amount"`
	Quantity      int    `bson:"quantity"`
	Currency      string `bson:"currency"`
	URL           string `bson:"url"`
}
