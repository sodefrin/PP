package lib

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	slog.SetDefault(logger)

	tests := []struct {
		name       string
		statusCode int
		wantLevel  string
	}{
		{"Success", 200, "INFO"},
		{"ClientError", 400, "INFO"},
		{"ServerError", 500, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			logOutput := buf.String()
			if !strings.Contains(logOutput, `"level":"`+tt.wantLevel+`"`) {
				t.Errorf("Expected log level %s, got log: %s", tt.wantLevel, logOutput)
			}
		})
	}
}
