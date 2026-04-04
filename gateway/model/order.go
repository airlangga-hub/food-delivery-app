package model

import "github.com/google/uuid"

type Order struct {
	OrderID           uuid.UUID
	RestaurantName    string
	DriverName        string
	ItemName          string
	RestaurantAddress string
	DeliveryAddress   string
}
