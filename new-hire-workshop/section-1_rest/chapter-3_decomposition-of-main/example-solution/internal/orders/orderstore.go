package orders

import "time"

type OrderStore struct {
	orders []Order
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: []Order{
			{
				CustomerID: "customer1",
				OrderID:    "order-1",
				Address:    "1 Longdown Avenue, BUK03",
				ProductID:  "HPE Alletra",
				PaymentID:  "1",
				Status:     OrderStatusPlaced,
				DeliveryEntries: []DeliveryEntry{
					{
						Timestamp: time.Now(),
						Message:   "Order has been placed",
					},
				},
				OrderedAt:   time.Now(),
				DeliveredAt: time.Time{},
			},
			{
				CustomerID: "customer2",
				OrderID:    "order-2",
				Address:    "Pepenero",
				ProductID:  "Margerita",
				PaymentID:  "2",
				Status:     OrderStatusPlaced,
				DeliveryEntries: []DeliveryEntry{
					{
						Timestamp: time.Now(),
						Message:   "Order has been placed",
					},
				},
				OrderedAt:   time.Now(),
				DeliveredAt: time.Time{},
			},
		},
	}
}

func (o *OrderStore) GetOrders(customerID string) []Order {
	var customerOrders []Order
	for _, order := range o.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	return customerOrders
}
