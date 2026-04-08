package model

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleDriver   UserRole = "driver"
)

type UserRegister struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Address     string
	PhoneNumber string
}

type UserInfo struct {
	FirstName   string
	LastName    string
	Email       string
	Address     string
	PhoneNumber string
	Balance     int
}

type PaymentLink struct {
	PaymentLink string
}
