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

func New(db *sql.DB) *sqlRepository {
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
