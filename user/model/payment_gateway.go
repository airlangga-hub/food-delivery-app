package model

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
	URL           string `bson:"url"`
}
