package routes

import (
	"log/slog"
	"net/http"
	"time"
)

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *logResponseWriter) WriteHeader(code int) {
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func LogIncomingRequest(r *http.Request) {
	slog.Info("Incoming Request...", "addr", r.RemoteAddr, "path", r.URL.Path, "method", r.Method)
}

func LogOutboundRequest(r *http.Request, statusCode int, duration time.Duration) {
	slog.Info("Received Request...", "addr", r.RemoteAddr, "path", r.URL.Path, "method", r.Method, "code", statusCode, "time", duration)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		LogIncomingRequest(r)
		// capture the response
		RW := &logResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// handle panics
		defer func() {
			r := recover()
			if r != nil {
				slog.Error("Panic occurred...", "context", r)
				RW.status = http.StatusInternalServerError
				http.Error(RW, "Internal server error", RW.status)
			}
		}()

		// send to next handler func
		next(RW, r.WithContext(ctx))

		duration := time.Since(start)
		LogOutboundRequest(r, RW.status, duration)
	}
}
