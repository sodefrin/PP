package api

import (
	"encoding/json"
	"net/http"

	"github.com/sodefrin/PP/server/api/dto"
	"github.com/sodefrin/PP/server/lib"
)

func MeHandler() lib.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return nil
		}

		user, err := lib.GetUserContext(r.Context())
		if err != nil {
			return err
		}

		resp := dto.User{
			ID:   user.ID,
			Name: user.Name,
		}

		respJSON, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(respJSON); err != nil {
			return err
		}
		return nil
	}
}
