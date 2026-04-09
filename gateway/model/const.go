package model

type (
	EmailStatus               string
	PaymentType               string
	LedgerReason              string
	OrderStatus               string
	PaymentGatewayRefIDPrefix string
	RoleUser                  string
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

	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"

	PaymentGatewayRefIDPrefixTopUp PaymentGatewayRefIDPrefix = "TOPUP_"
	PaymentGatewayRefIDPrefixOrder PaymentGatewayRefIDPrefix = "ORDER_"

	RoleUserCustomer RoleUser = "customer"
	RoleUserDriver   RoleUser = "driver"
)
