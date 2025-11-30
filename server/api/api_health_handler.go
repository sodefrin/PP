package api

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
