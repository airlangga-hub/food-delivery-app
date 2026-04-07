package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"
)

type Order struct {
	ID                  uuid.UUID
	Restaurants         []Restaurant
	DeliveryAddress     string
	CustomerPhoneNumber string
	OrderStatus         OrderStatus
	Driver              Driver
	DeliveryFee         int
	TotalFee            int
}

type Restaurant struct {
	ID      uuid.UUID
	Name    string
	Address string
	Items   []Item
}

type Item struct {
	ID       uuid.UUID
	Name     string
	Price    int
	Quantity int
}

type Driver struct {
	ID            uuid.UUID
	AverageRating float64
	Name          string
	Bike          string
	LicensePlate  string
	PhoneNumber   string
}

type FindDriver struct {
	OrderApplicants []Driver
}
