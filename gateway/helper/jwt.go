package helper

import (
	"github.com/airlangga-hub/food-delivery-app/gateway/model"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID string
	Email  string
	Role   model.RoleUser
}
