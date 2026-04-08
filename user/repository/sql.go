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

func (r *sqlRepository) RegisterCustomer(ctx context.Context, input model.UserRegister) (model.UserInfo, error) {
	_, err := r.db.ExecContext(
		ctx,
		`WITH new_user AS (
			INSERT INTO
				users (email, password_hash, role)
			VALUES
				($1, $2, $7)
			RETURNING
				id
		)
		INSERT INTO
			customer_profiles (user_id, first_name, last_name, address, phone_number)
		SELECT
			id, $3, $4, $5, $6
		FROM
			new_user
		RETURNING
			user_id`,
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
		input.Address,
		input.PhoneNumber,
		model.UserRoleCustomer,
	)

	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.repository.RegisterCustomer (ExecContext): %w", err)
	}

	return model.UserInfo{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
	}, nil
}

func (r *sqlRepository) Login(ctx context.Context, email string) (string, error) {
	var passwordHash string

	err := r.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE email = $1`, email).Scan(&passwordHash)
	if err != nil {
		return "", fmt.Errorf("user.repository.Login (QueryRowContext): %w", err)
	}

	return passwordHash, nil
}

func (r *sqlRepository) GetUserInfo(ctx context.Context, email string) (model.UserInfo, error) {
	var u model.UserInfo
	query := `SELECT first_name, last_name, email, address, balance FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &u.Address, &u.Balance)
	return u, err
}
