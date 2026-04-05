package handler

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
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

type TopUpBalanceParam struct {
	Amount int `query:"amount" validate:"required"`
}

type CreateOrderParam struct {
	ItemID string `query:"item_id" validate:"required"`
}

type GetDriversParam struct {
	OrderID string `query:"order_id" validate:"required"`
}