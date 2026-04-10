package handler

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TopUpBalanceRequest struct {
	Amount int `json:"amount" validate:"required"`
}

type ItemRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CreateOrderRequest struct {
	DeliveryFee int           `json:"delivery_fee"`
	Items       []ItemRequest `json:"items" validate:"required"`
}

type GiveRatingRequest struct {
	Rating int `json:"rating" validate:"required"`
}

type ChooseDriverRequest struct {
	DriverID string `json:"driver_id" validate:"required"`
}

type XenditWebhookRequest struct {
	Created    string `json:"created"`    
	BusinessID string `json:"business_id"`
	Event      string `json:"event"`      
	APIVersion string `json:"api_version"`
	Data       Data   `json:"data" validate:"required"`       
}

type Data struct {
	PaymentSessionID string `json:"payment_session_id"`
	Created          string `json:"created"`           
	Updated          string `json:"updated"`           
	Status           string `json:"status"`            
	ReferenceID      string `json:"reference_id" validate:"required"`      
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
	PaymentRequestID string `json:"payment_request_id"`
	PaymentID        string `json:"payment_id"`        
	PaymentLinkURL   string `json:"payment_link_url"`  
}

type Item struct {
	ReferenceID   string `json:"reference_id"`   
	Type          string `json:"type"`           
	Name          string `json:"name"`           
	NetUnitAmount int64  `json:"net_unit_amount"`
	Quantity      int64  `json:"quantity"`       
	URL           string `json:"url"`            
	Category      string `json:"category"`       
	Description   string `json:"description"`    
}