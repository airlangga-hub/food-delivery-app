package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/airlangga-hub/food-delivery-app/user/model"
	"github.com/google/uuid"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *sqlRepository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) UpdateLedger(ctx context.Context, userID uuid.UUID, reason model.LedgerReason, amount int) error {
	finalAmount := amount

	if reason == model.LedgerReasonCustomerOrder {
		finalAmount = -amount
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO
			ledgers (user_id, amount, reason)
		VALUES
			($1, $2, $3)`,
		userID, finalAmount, string(reason),
	)
	if err != nil {
		return fmt.Errorf("user.repository.UpdateLedger: %w", err)
	}

	return nil
}

func (r *sqlRepository) Register(ctx context.Context, u model.UserRegister) error {
	query := `INSERT INTO users (first_name, last_name, email, password_hash, address) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Address)
	return err
}

func (r *sqlRepository) Login(ctx context.Context, email string) (model.UserInfo, string, error) {
	var u model.UserInfo
	var password string
	query := `SELECT first_name, last_name, email, password_hash, address, balance FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.FirstName, &u.LastName, &u.Email, &u.Address, &u.Balance, &password,
	)
	return u, password, err
}

func (r *sqlRepository) GetUserInfo(ctx context.Context, email string) (model.UserInfo, error) {
	var u model.UserInfo
	query := `SELECT first_name, last_name, email, address, balance FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &u.Address, &u.Balance)
	return u, err
}
