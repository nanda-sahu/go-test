type RESTOrderStatus string

const (
	RESTOrderStatusPlaced    RESTOrderStatus = "Placed"
	RESTOrderStatusDelivered RESTOrderStatus = "Delivered"
	RESTOrderStatusTransit   RESTOrderStatus = "In transit"
)

type ListOrdersResponse struct {
	OrderSummaries []OrderSummary `json:"orderSummaries"`
}

type OrderSummary struct {
	OrderID    string          `json:"orderId"`
	CustomerID string          `json:"customerId"`
	ProductID  string          `json:"productId"`
	Status     RESTOrderStatus `json:"status"`
}

type Error struct {
	Message string `json:"message"`
}

