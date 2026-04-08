package model

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

	LedgerReasonTopUp               LedgerReason = "customer_top_up"
	LedgerReasonCustomerOrder       LedgerReason = "customer_order"
	LedgerReasonDriverCompleteOrder LedgerReason = "driver_complete_order"
)

type PaymentRecord struct {
	ID                     string
	Email                  string
	EmailStatus            EmailStatus
	PaymentType            PaymentType
	PaymentGatewayResponse PaymentGatewayResponse
}
