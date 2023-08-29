package orders

import "time"

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "PLACED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusTransit   OrderStatus = "IN TRANSIT"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusOnHold    OrderStatus = "ON HOLD"
)

type DeliveryEntry struct {
	Timestamp time.Time
	Message   string
}

type Order struct {
	CustomerID      string
	OrderID         string
	Address         string
	ProductID       string
	PaymentID       string
	Status          OrderStatus
	DeliveryEntries []DeliveryEntry
	OrderedAt       time.Time
	DeliveredAt     time.Time
}
