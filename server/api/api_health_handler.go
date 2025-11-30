package api

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		// Can't do much if writing response fails, but good to know
		// slog.Error("Failed to write health response", "error", err)
	}
}
