package model

import "github.com/google/uuid"

type Order struct {
	OrderID           uuid.UUID `json:"order_id"`
	RestaurantName    string    `json:"restaurant_name"`
	DriverName        string    `json:"driver_name"`
	ItemName          string    `json:"item_name"`
	RestaurantAddress string    `json:"restaurant_address"`
	DeliveryAddress   string    `json:"delivery_address"`
}
