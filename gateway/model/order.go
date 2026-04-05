package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"
)

type Order struct {
	OrderID         uuid.UUID    `json:"order_id"`
	Restaurants     []Restaurant `json:"restaurants"`
	DeliveryAddress string       `json:"delivery_address"`
	OrderStatus     OrderStatus  `json:"order_status"`
	Driver          *Driver      `json:"driver,omitempty"`
}

type Restaurant struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Items   []Item `json:"items"`
}

type Driver struct {
	DriverID            uuid.UUID `json:"driver_id"`
	DriverAverageRating float64   `json:"driver_average_rating"`
	DriverName          string    `json:"driver_name"`
	DriverBike          string    `json:"driver_bike"`
	DriverLicensePlate  string    `json:"driver_license_plate"`
}

type FindDriver struct {
	OrderID          uuid.UUID `json:"order_id"`
	ApplicantDrivers []Driver  `json:"applicant_driver"`
}

type Item struct {
	ItemID   uuid.UUID
	ItemName string
	Quantity int
}
