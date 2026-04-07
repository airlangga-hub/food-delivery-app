package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"
)

type Order struct {
	ID                  uuid.UUID    `json:"order_id"`
	Restaurants         []Restaurant `json:"restaurants"`
	DeliveryAddress     string       `json:"delivery_address"`
	CustomerPhoneNumber string       `json:"customer_phone_number"`
	OrderStatus         OrderStatus  `json:"order_status"`
	Driver              *Driver      `json:"driver,omitempty"`
}

type Restaurant struct {
	ID      uuid.UUID `json:"restaurant_id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Items   []Item    `json:"items"`
}

type Item struct {
	ID       uuid.UUID `json:"item_id"`
	Name     string    `json:"item_name"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
}

type Driver struct {
	ID                  uuid.UUID `json:"driver_id"`
	DriverAverageRating float64   `json:"driver_average_rating"`
	DriverName          string    `json:"driver_name"`
	DriverBike          string    `json:"driver_bike"`
	DriverLicensePlate  string    `json:"driver_license_plate"`
}

type FindDriver struct {
	OrderID          uuid.UUID `json:"order_id"`
	ApplicantDrivers []Driver  `json:"applicant_driver"`
}
