package rest

import (
	"ecommerce-workshop/internal/orders"
)

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "Placed"
	OrderStatusDelivered OrderStatus = "Delivered"
	OrderStatusTransit   OrderStatus = "In transit"
)

type ListOrdersResponse struct {
	OrderSummaries []OrderSummary `json:"orderSummaries"`
}

type OrderSummary struct {
	OrderID    string      `json:"orderId"`
	CustomerID string      `json:"customerId"`
	ProductID  string      `json:"productId"`
	Status     OrderStatus `json:"status"`
}

type Error struct {
	Message string `json:"message"`
}

func internalOrdersToSummaries(orders []orders.Order) (ListOrdersResponse, error) {
	// Convert core Orders to the REST representation
}

func orderStatusToREST(status orders.OrderStatus) (OrderStatus, error) {
	// Convert core enum to REST value
}
