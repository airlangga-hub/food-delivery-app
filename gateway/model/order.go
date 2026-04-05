package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderLookingForDriver OrderStatus = "Looking for driver"
	OrderDriverOTW        OrderStatus = "Driver is on the way"
	OrderDone             OrderStatus = "Done"
)

type Order struct {
	OrderID           uuid.UUID   `json:"order_id"`
	RestaurantName    string      `json:"restaurant_name"`
	RestaurantAddress string      `json:"restaurant_address"`
	ItemName          string      `json:"item_name"`
	DeliveryAddress   string      `json:"delivery_address"`
	OrderStatus       OrderStatus `json:"order_status"`
	Driver            *Driver     `json:"driver,omitempty"`
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
