package model

type PaymentRecord struct {
	ID                     string
	Email                  string
	EmailStatus            EmailStatus
	PaymentType            PaymentType
	PaymentGatewayResponse PaymentGatewayResponse
}
