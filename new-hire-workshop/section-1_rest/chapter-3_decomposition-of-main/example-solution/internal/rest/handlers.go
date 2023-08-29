package rest

import (
	"ecommerce-workshop/internal/orders"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.hpe.com/cloud/go-gadgets/x/logging"
)

type ListOrderSummariesHandler struct {
	orderStore *orders.OrderStore
	logger     logging.Logger
}

func NewListOrderSummariesHandler(orderStore *orders.OrderStore, logger logging.Logger) *ListOrderSummariesHandler {
	return &ListOrderSummariesHandler{
		orderStore: orderStore,
		logger:     logger,
	}
}

func (l *ListOrderSummariesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customerID")
	if customerID == "" {
		writeError(w, "customerID not provided", http.StatusBadRequest, l.logger)
		return
	}

	orders := l.orderStore.GetOrders(customerID)

	summary, err := internalOrdersToSummaries(orders)
	if err != nil {
		l.logger.WithError(err).WithField("orders", orders).Error("failed to convert orders from core to rest")
		writeError(w, "An internal error occurred", http.StatusInternalServerError, l.logger)
		return
	}

	respBody, err := json.Marshal(summary)
	if err != nil {
		l.logger.WithError(err).WithField("summary", summary).Error("failed to marshal order summaries")
		writeError(w, "An internal error occurred", http.StatusInternalServerError, l.logger)
		return
	}

	if _, err := w.Write(respBody); err != nil {
		l.logger.WithError(err).WithField("resp-body", respBody).Error("failed to write response body")
		writeError(w, "An internal error occurred", http.StatusInternalServerError, l.logger)
		return
	}
}

func writeError(w http.ResponseWriter, msg string, status int, logger logging.Logger) {
	restErr := Error{
		Message: msg,
	}

	errResp, err := json.Marshal(restErr)
	if err != nil {
		logger.WithError(err).WithField("rest-err", restErr).Error("failed to marshall error response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(errResp); err != nil {
		logger.WithError(err).WithField("err-resp", errResp).Error("failed to write error response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
