package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sodefrin/PP/server/db"
	"github.com/sodefrin/PP/server/lib"
)

func MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user := r.Context().Value(lib.UserContextKey).(db.User)

		// Don't return password hash
		resp := struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		}

		respJSON, err := json.Marshal(resp)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(respJSON); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	}
}
