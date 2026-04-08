package repository

import (
	"context"
	"database/sql"

	"github.com/airlangga-hub/food-delivery-app/user/model"
)

type UserRepository interface {
	Register(ctx context.Context, user model.UserRegister) error
	Login(ctx context.Context, email string) (model.UserInfo, string, error)
	GetUserByInfo(ctx context.Context, email string) (model.UserInfo, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(ctx context.Context, u model.UserRegister) error {
	query := `INSERT INTO users (first_name, last_name, email, password, address) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Address)
	return err
}

func (r *userRepository) Login(ctx context.Context, email string) (model.UserInfo, string, error) {
	var u model.UserInfo
	var password string
	query := `SELECT first_name, last_name, email, password, address, balance FROM users WHERE = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.FirstName, &u.LastName, &u.Email, &u.Address, &u.Balance, &password,
	)
	return u, password, err
}

func (r *userRepository) GetUserByInfo(ctx context.Context, email string) (model.UserInfo, error) {
	var u model.UserInfo
	query := `SELECT first_name, last_name, email, address, balance FROM users WHERE = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &u.Address, &u.Balance)
	return u, err
}
