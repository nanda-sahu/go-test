package rest

import (
	"ecommerce-workshop/internal/orders"
	"net/http"

	"github.hpe.com/cloud/go-gadgets/x/logging"
)

type ListOrderSummariesHandler struct {
	orderStore *orders.OrderStore
	logger     logging.Logger
}

func NewListOrderSummariesHandler(orderStore *orders.OrderStore, logger logging.Logger) *ListOrderSummariesHandler {
	// Populate and return handler here
}

func (l *ListOrderSummariesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Migrate across the handler logic into this function
}
