package service

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *userService {
	return &userService{db: db}
}

