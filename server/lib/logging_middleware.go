package lib

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Extract trace ID
		span := trace.SpanFromContext(r.Context())
		traceID := span.SpanContext().TraceID().String()

		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", duration),
			slog.String("trace_id", traceID),
		}

		if rw.statusCode >= 500 {
			slog.Error("HTTP Request Failed", attrs...)
		} else {
			slog.Info("HTTP Request", attrs...)
		}
	})
}
