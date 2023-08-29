var orders = []*Order{
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
}
