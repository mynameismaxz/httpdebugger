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

// DebugHandler echoes the details of any incoming request back to the client.
// It renders an HTML page by default, or JSON when the client asks for it via
// the Accept header or a ?format=json query parameter.
func (h *HttpRouteHandler) DebugHandler(w http.ResponseWriter, r *http.Request) {
	info := captureRequest(w, r)

	h.logger.Info("request captured",
		"method", info.Method,
		"path", info.Path,
		"client_ip", info.ClientIP,
		"proto", info.Proto,
		"content_type", info.ContentType,
		"body_size", info.BodySize,
	)

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(info.ToJSON())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := debugTmpl.Execute(w, newDebugView(info)); err != nil {
		h.logger.Error("failed to render debug template: " + err.Error())
	}
}

// HealthHandler reports server health status as JSON. It is kept separate from
// the debugger so platforms like Cloud Run can probe it at a stable path.
func (h *HttpRouteHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
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
