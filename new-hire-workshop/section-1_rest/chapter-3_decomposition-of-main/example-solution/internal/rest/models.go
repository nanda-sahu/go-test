package rest

import (
	"ecommerce-workshop/internal/orders"
	"fmt"
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
	resp := ListOrdersResponse{
		OrderSummaries: make([]OrderSummary, 0, len(orders)),
	}

	for _, order := range orders {
		status, err := orderStatusToREST(order.Status)
		if err != nil {
			return ListOrdersResponse{}, err
		}

		summary := OrderSummary{
			OrderID:    order.OrderID,
			CustomerID: order.CustomerID,
			ProductID:  order.ProductID,
			Status:     status,
		}

		resp.OrderSummaries = append(resp.OrderSummaries, summary)
	}

	return resp, nil
}

func orderStatusToREST(status orders.OrderStatus) (OrderStatus, error) {
	switch status {
	case orders.OrderStatusPlaced, orders.OrderStatusCancelled, orders.OrderStatusOnHold:
		return OrderStatusPlaced, nil
	case orders.OrderStatusDelivered:
		return OrderStatusDelivered, nil
	case orders.OrderStatusTransit:
		return OrderStatusTransit, nil
	}

	return OrderStatus(""), fmt.Errorf("Unexpected status value: %s", status)
}
