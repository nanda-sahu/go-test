package main

import (
	"context"
	"example-solution/internal/orders"
	"example-solution/internal/rest"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.hpe.com/cloud/go-gadgets/x/logging"
)

const (
	ErrLoggerInit         = 1
	ErrServerInit         = 2
	ErrShuttingDownServer = 3
)

func run() int {
	logger, err := logging.NewZapJSONLogger("info")
	if err != nil {
		// If we fail to init the logger, we log out the error in a non-structured manner as it's
		// better to have the error message, even if it's not in the right format
		fmt.Printf("failed to init logger: %s", err.Error())
		return ErrLoggerInit
	}

	baseCtx := context.Background()
	server, err := makeServer(baseCtx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to create server")
		return ErrServerInit
	}

	ctx, stop := signal.NotifyContext(baseCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Call `server.Start()` inside of a goroutine so it does not block the goroutine that the application
	// is running on. This is mainly if we're running multiple servers/processes (e.g. a REST server and
	// a metrics server)
	go server.Start()

	<-ctx.Done()

	if err := server.Stop(); err != nil {
		logger.WithError(err).Error("error shutting down rest server")
		return ErrShuttingDownServer
	}

	return 0
}

func makeServer(baseCtx context.Context, logger logging.Logger) (*rest.Server, error) {
	orderStore := orders.NewOrderStore()

	router := rest.NewMux(orderStore, logger)

	return rest.NewServer(baseCtx, 8080, router, logger)
}

// It looks quite strange to have main be such a small bit of code for main, but we
// want to avoid having many places in the code where we call os.Exit (which we do to ensure
// we get a non-zero return code for the application on error).
func main() {
	if retCode := run(); retCode != 0 {
		os.Exit(retCode)
	}
}
