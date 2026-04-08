package repository

import (
	"context"
	"database/sql"

	"github.com/airlangga-hub/food-delivery-app/user/model"
)

type UserRepository interface {
	Register(ctx context.Context, user model.UserRegister) error
	Login(ctx context.Context, email string) (model.UserInfo, string, error)
	FindUserByEmail(ctx context.Context, email string) (model.UserInfo, string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(ctx context.Context, user model.UserRegister) error {

}
