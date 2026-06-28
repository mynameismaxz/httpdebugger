package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mynameismaxz/httpdebugger/handlers"
)

const (
	// DefaultHttpServerPort is used when the PORT environment variable is unset.
	// Cloud Run injects PORT (default 8080) and expects the server to listen on it.
	DefaultHttpServerPort = "8080"

	// shutdownTimeout bounds how long in-flight requests have to drain on SIGTERM.
	shutdownTimeout = 10 * time.Second
)

func main() {
	// Create logger for instance server
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Resolve the listen port from the environment (Cloud Run sets PORT).
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultHttpServerPort
	}

	// Create Handler for HTTP Route
	httpHandler := handlers.New(logger)

	// Use a dedicated mux so handlers are scoped to this server instance
	// rather than the global default mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", httpHandler.HealthHandler)
	mux.HandleFunc("/", httpHandler.DebugHandler) // catch-all: any path other than /healthz

	// Register server configuration
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	// Start the server in a goroutine so main can wait for shutdown signals.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("Server started at port " + port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Wait for either a fatal server error or a termination signal.
	// Cloud Run sends SIGTERM before stopping the container instance.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Error("Server error: " + err.Error())
		os.Exit(1)
	case sig := <-stop:
		logger.Info("Shutdown signal received: " + sig.String())
	}

	// Gracefully drain in-flight requests within the timeout window.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed: " + err.Error())
		os.Exit(1)
	}

	logger.Info("Server stopped cleanly")
}
