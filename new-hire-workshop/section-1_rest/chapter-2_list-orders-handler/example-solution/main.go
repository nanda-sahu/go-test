package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

var orders = []Order{
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/orders", ListOrderSummariesHandler).Methods(http.MethodGet)

	port := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Broadcast the HTTP server on port 8080 of localhost, with an handler on the path "/".
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

// Create a handler function that responds to a request with a string saying "Hello, "/"".
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello on path %q", html.EscapeString(r.URL.Path))
}

func ListOrderSummariesHandler(w http.ResponseWriter, r *http.Request) {
	// Whilst we can use mux to get path parameters, to get query parameters
	// we need to use the net/http stdlib functionality
	customerID := r.URL.Query().Get("customerID")
	if customerID == "" {
		writeError(w, "customerID not provided", http.StatusBadRequest)
		return
	}

	orders := getCustomerOrders(customerID)

	summary, err := internalOrdersToRESTSummaries(orders)
	if err != nil {
		// In a production situation, we'd want to log the error for ourselves
		// along with some tracing info so that we could debug it.
		// We don't want to return internal error messages to the consumer.
		writeError(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(summary)
	if err != nil {
		writeError(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respBody); err != nil {
		writeError(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}
}

// writeError is a simple utility function for error responses, used to keep the handler code
// cleaner and avoid duplication.
func writeError(w http.ResponseWriter, msg string, status int) {
	restErr := Error{
		Message: msg,
	}

	errResp, err := json.Marshal(restErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(errResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getCustomerOrders abstracts away fetching the customer's orders. At this point it's just a simple
// case of looping through our slice and returning any that match the specified ID.
func getCustomerOrders(customerID string) []Order {
	var customerOrders []Order
	for _, order := range orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	return customerOrders
}

// internalOrderToRESTSummaries handles converting the internal data model to the external
// REST representation to be used as the response body.
func internalOrdersToRESTSummaries(orders []Order) (ListOrdersResponse, error) {
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

// orderStatusToREST converts our internal enum to the REST enum. By having the error handling
// at the end of the function rather than a 'default' case the linter can tell us if we've missed
// any possible values.
func orderStatusToREST(status OrderStatus) (RESTOrderStatus, error) {
	switch status {
	case OrderStatusPlaced, OrderStatusCancelled, OrderStatusOnHold:
		return RESTOrderStatusPlaced, nil
	case OrderStatusDelivered:
		return RESTOrderStatusDelivered, nil
	case OrderStatusTransit:
		return RESTOrderStatusTransit, nil
	}

	return RESTOrderStatus(""), fmt.Errorf("Unexpected status value: %s", status)
}
