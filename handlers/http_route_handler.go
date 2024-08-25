package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/mynameismaxz/httpdebugger/dtos"
)

type HttpRouteHandler struct {
	logger *slog.Logger
}

func New(l *slog.Logger) *HttpRouteHandler {
	return &HttpRouteHandler{
		logger: l,
	}
}

func (h *HttpRouteHandler) MainRouteHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// get hostname
	hostName, err := os.Hostname()
	if err != nil {
		hostName = ""
	}

	// count runtime
	runtimeCount := runtime.NumGoroutine()

	// get log level

	response := &dtos.MainRouteResponse{
		Status:  http.StatusOK,
		Message: "Server is running",
		Server: dtos.ServerDetailStatus{
			Timestamp:      time.Now().UTC().Format(time.RFC3339),
			ProcessingTime: time.Since(startTime).String(),
			HostName:       hostName,
		},
		Application: dtos.ApplicationDetailStatus{
			LogLevel:       slog.LevelInfo.String(),
			GoRoutineCount: runtimeCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response.ToJSON())
}
