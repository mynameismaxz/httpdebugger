package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/mynameismaxz/httpdebugger/handlers"
)

const (
	// DefaultPort for the server
	DefaultHttpServerPort = "1337"
)

func main() {
	// Create logger for instance server
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Create Handler for HTTP Route
	httpHandler := handlers.New(logger)

	// Register server configuration
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", DefaultHttpServerPort),
	}

	// Register handler for route
	http.HandleFunc("/", httpHandler.MainRouteHandler)

	// Start server
	logger.Info("Server started at port " + DefaultHttpServerPort)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server error: " + err.Error())
	}
}
