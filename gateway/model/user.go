package model

type UserRegister struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone_number"`
}

type UserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address   string `json:"address"`
	Balance   int    `json:"balance"`
}

type PaymentLink struct {
	PaymentLink string `json:"payment_link"`
}
