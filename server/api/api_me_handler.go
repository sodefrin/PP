package api

import (
	"encoding/json"
	"net/http"

	"github.com/sodefrin/PP/server/db"
	"github.com/sodefrin/PP/server/lib"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value(lib.UserContextKey).(db.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Don't return password hash
	resp := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   user.ID,
		Name: user.Name,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
