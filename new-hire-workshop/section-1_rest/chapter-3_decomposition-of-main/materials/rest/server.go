package rest

import (
	"context"
	"ecommerce-workshop/internal/orders"
	"net/http"

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
	// Do some basic input validation

	// Populate and return server
}

func (s *Server) Start() {
	// Start the server
}

func (s *Server) Stop() error {
	// Stop the server, with a timeout for stopping it
}

func NewMux(orderStore *orders.OrderStore, logger logging.Logger) *mux.Router {
	// Create a mux and register our handlers
}
