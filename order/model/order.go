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
	ID                  uuid.UUID
	DriverAverageRating float64
	DriverName          string
	DriverBike          string
	DriverLicensePlate  string
}

type FindDriver struct {
	OrderID          uuid.UUID
	ApplicantDrivers []Driver
}
