package repository

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"
)

type Restaurant struct {
	ID        uuid.UUID
	Name      string
	Address   string
	CreatedAt time.Time
}

type Item struct {
	ID        uuid.UUID
	Name      string
	Stock     int
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID        uuid.UUID
	Name      string
	Quantity  int
	CreatedAt time.Time
}

type Order struct {
	ID          uuid.UUID
	OrderStatus OrderStatus
	DeliveryFee int
	TotalFee    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Customer struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Driver struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Bike         string
	LicensePlate string
	CreatedAt    time.Time
}

type Rating struct {
	ID        uuid.UUID
	Rating    int
	CreatedAt time.Time
}

type Ledger struct {
	ID        uuid.UUID
	Amount    int
	Reason    string
	CreatedAt time.Time
}
