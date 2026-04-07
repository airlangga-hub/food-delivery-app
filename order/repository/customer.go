package repository

import "database/sql"

type customerRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *customerRepository {
	return &customerRepository{db: db}
}