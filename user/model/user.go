package model

import "github.com/google/uuid"

type UserRegister struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Address     string
	PhoneNumber string
}

type UserLogin struct {
	UserID       uuid.UUID
	PasswordHash string
	Email        string
	Role         string
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
