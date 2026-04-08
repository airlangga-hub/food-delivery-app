package model

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
	Email                  string
	EmailStatus            EmailStatus
	PaymentType            PaymentType
	PaymentGatewayResponse PaymentGatewayResponse
}
