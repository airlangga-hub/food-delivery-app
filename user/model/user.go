package model

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleDriver   UserRole = "driver"
)

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Address   string `json:"address"`
}

type UserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Balance   int    `json:"balance"`
}

type PaymentLink struct {
	PaymentLink string `json:"payment_link"`
}
