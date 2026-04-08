package repository

import (
	"context"
	"database/sql"

	"github.com/airlangga-hub/food-delivery-app/user/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.UserRegister) error
	FindUserByEmail(ctx context.Context, email string) (model.UserInfo, string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
