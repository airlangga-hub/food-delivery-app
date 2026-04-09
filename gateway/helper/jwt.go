package helper

import (
	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleCustomer Role = "customer"
	RoleDriver   Role = "driver"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID string
	Email  string
	Role   Role
}
