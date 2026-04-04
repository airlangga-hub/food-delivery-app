package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID int
	Email  string
	Role   string
}

func MakeJWT(userID int, email, role string, key []byte) (string, error) {
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
