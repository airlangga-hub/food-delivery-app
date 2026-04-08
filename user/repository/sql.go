package repository

import (
	"database/sql"
)

type sqlRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *sqlRepository {
	return &sqlRepository{db: db}
}
