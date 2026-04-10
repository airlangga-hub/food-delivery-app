package repository

import (
	"context"
	"database/sql"
	"errors"
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
			final_project.ledgers (user_id, amount, reason)
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
				final_project.users (email, password_hash, role)
			VALUES
				($1, $2, $7)
			RETURNING
				id
		)
		INSERT INTO
			final_project.customer_profiles (user_id, first_name, last_name, address, phone_number)
		SELECT
			id, $3, $4, $5, $6
		FROM
			new_user`,
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
		input.Address,
		input.PhoneNumber,
		model.RoleUserCustomer,
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

func (r *sqlRepository) Login(ctx context.Context, email string) (model.UserLogin, error) {
	var user model.UserLogin

	err := r.db.QueryRowContext(ctx, `SELECT id, password_hash, email, role FROM final_project.users WHERE email = $1`, email).Scan(
		&user.UserID,
		&user.PasswordHash,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserLogin{}, fmt.Errorf("user.repository.Login (no rows found): %w", model.ErrNotFound)
		}
		return model.UserLogin{}, fmt.Errorf("user.repository.Login (Scan): %w", err)
	}

	return user, nil
}

func (r *sqlRepository) GetUserInfo(ctx context.Context, email string, role model.RoleUser) (model.UserInfo, error) {
	var u model.UserInfo

	var tableName string
	var addressCol string

	switch role {
	case model.RoleUserCustomer:
		tableName = "customer_profiles"
		addressCol = "p.address"
	case model.RoleUserDriver:
		tableName = "driver_profiles"
		addressCol = "''"
	default:
		return model.UserInfo{}, fmt.Errorf("invalid role: %s", role)
	}

	err := r.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`
				SELECT
					p.first_name,
					p.last_name,
					u.email,
					%s AS address,
					p.phone_number,
					COALESCE((
						SELECT SUM(amount)
						FROM final_project.ledgers
						WHERE user_id = u.id
					), 0) AS balance
				FROM
					final_project.users u
				JOIN
					final_project.%s p ON u.id = p.user_id
				WHERE
					u.email = $1`,
			addressCol,
			tableName,
		),
		email,
	).Scan(
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Address,
		&u.PhoneNumber,
		&u.Balance,
	)

	if err != nil {
		return model.UserInfo{}, fmt.Errorf("user.repository.GetUserInfo (Scan): %w", err)
	}

	return u, nil
}
