package repository

type XenditPaymentSessionRequest struct {
	ReferenceID      string   `json:"reference_id" validate:"required"`
	SessionType      string   `json:"session_type" validate:"required"`
	Mode             string   `json:"mode" validate:"required"`
	Amount           int      `json:"amount" validate:"required"`
	Currency         string   `json:"currency" validate:"required"`
	Country          string   `json:"country" validate:"required"`
	Customer         Customer `json:"customer" validate:"required"`
	Items            []Item   `json:"items,omitempty"`
	CaptureMethod    string   `json:"capture_method,omitempty"`
	Locale           string   `json:"locale,omitempty"`
	Description      string   `json:"description,omitempty"`
	SuccessReturnURL string   `json:"success_return_url,omitempty"`
	CancelReturnURL  string   `json:"cancel_return_url,omitempty"`
}

type XenditPaymentSessionResponse struct {
	PaymentSessionID string `json:"payment_session_id"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
	Status           string `json:"status"`
	ReferenceID      string `json:"reference_id"`
	Currency         string `json:"currency"`
	Amount           int64  `json:"amount"`
	Country          string `json:"country"`
	ExpiresAt        string `json:"expires_at"`
	SessionType      string `json:"session_type"`
	Mode             string `json:"mode"`
	Locale           string `json:"locale"`
	BusinessID       string `json:"business_id"`
	CustomerID       string `json:"customer_id"`
	CaptureMethod    string `json:"capture_method"`
	Description      string `json:"description"`
	Items            []Item `json:"items"`
	SuccessReturnURL string `json:"success_return_url"`
	CancelReturnURL  string `json:"cancel_return_url"`
	PaymentLinkURL   string `json:"payment_link_url"`
}

type Customer struct {
	ReferenceID      string           `json:"reference_id" validate:"required"`
	Type             string           `json:"type" validate:"required"`
	Email            string           `json:"email,omitempty"`
	MobileNumber     string           `json:"mobile_number,omitempty"`
	IndividualDetail IndividualDetail `json:"individual_detail" validate:"required"`
}

type IndividualDetail struct {
	GivenNames string `json:"given_names" validate:"required"`
	Surname    string `json:"surname,omitempty"`
}

type Item struct {
	ReferenceID   string `json:"reference_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description,omitempty"`
	Type          string `json:"type,omitempty"`
	Category      string `json:"category" validate:"required"`
	NetUnitAmount int    `json:"net_unit_amount" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required"`
	Currency      string `json:"currency,omitempty"`
	URL           string `json:"url,omitempty"`
}
