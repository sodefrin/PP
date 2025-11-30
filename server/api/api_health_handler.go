package api

import (
	"net/http"

	"github.com/sodefrin/PP/server/lib"
)

func HealthHandler() lib.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if _, err := w.Write([]byte("OK")); err != nil {
			return err
		}
		return nil
	}
}
