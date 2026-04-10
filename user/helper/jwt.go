package helper

import (
	"time"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID string
	Email  string
	Role   model.RoleUser
}

func MakeJWT(userID string, email string, role model.RoleUser, key []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&MyClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   email,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
			UserID: userID,
			Role:   role,
		},
	)

	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
