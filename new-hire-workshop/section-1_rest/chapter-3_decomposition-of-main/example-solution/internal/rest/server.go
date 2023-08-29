package rest

import (
	"context"
	"ecommerce-workshop/internal/orders"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.hpe.com/cloud/go-gadgets/x/logging"
)

type Server struct {
	server *http.Server
	logger logging.Logger
}

// The baseCtx that is specified here is used when creating a context for all requests that this server
// creates. This allows us to cancel all currently handled requests from where we are invoking the server,
// should we wish to.
func NewServer(baseCtx context.Context, port int, router *mux.Router, logger logging.Logger) (*Server, error) {
	if router == nil {
		return nil, errors.New("router is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if baseCtx == nil {
		return nil, errors.New("baseCtx is nil")
	}

	return &Server{
		server: &http.Server{
			Handler: router,
			Addr:    fmt.Sprintf(":%d", port),
			BaseContext: func(net.Listener) context.Context {
				return baseCtx
			},
		},
		logger: logger,
	}, nil
}

func (s *Server) Start() {
	// We log out the address that the server is using for visibility.
	s.logger.WithField("address", s.server.Addr).Info("serving REST")

	// Broadcast the HTTP server that we created earlier, and if there is an
	// error that is not "the server is closed", we log it to stdout.
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.WithError(err).Error("non close error from server")
	}
}

func (s *Server) Stop() error {
	// Log out that the REST server has stopped for visibility.
	s.logger.WithField("address", s.server.Addr).Info("Shutting down REST server")

	// Create a context that will shut down the server in 5 seconds, no matter
	// what. We then pass this context to the server so it can do that if there
	// are any outstanding connections.
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutDownCtx)
}

func NewMux(orderStore *orders.OrderStore, logger logging.Logger) *mux.Router {
	listHandler := NewListOrderSummariesHandler(orderStore, logger)

	router := mux.NewRouter()
	router.Handle("api/v1/orders", listHandler).Methods(http.MethodGet)
	return router
}
