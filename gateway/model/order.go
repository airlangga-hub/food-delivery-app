package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusSearchingForDriver OrderStatus = "pending"
	OrderStatusDriverOTW          OrderStatus = "otw"
	OrderStatusDone               OrderStatus = "done"
)

type Order struct {
	ID                  uuid.UUID    `json:"id"`
	Restaurants         []Restaurant `json:"restaurants"`
	DeliveryAddress     string       `json:"delivery_address"`
	CustomerPhoneNumber string       `json:"customer_phone_number"`
	OrderStatus         OrderStatus  `json:"order_status"`
	Driver              *Driver      `json:"driver,omitempty"`
}

type Restaurant struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Items   []Item    `json:"items"`
}

type Item struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"item_name"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
}

type Driver struct {
	ID            uuid.UUID `json:"id"`
	AverageRating float64   `json:"average_rating"`
	Name          string    `json:"name"`
	Bike          string    `json:"bike"`
	LicensePlate  string    `json:"license_plate"`
	PhoneNumber   string    `json:"phone_number"`
}

type FindDriver struct {
	OrderApplicants []Driver `json:"order_applicants"`
}
