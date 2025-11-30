package lib

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sodefrin/PP/server/db"
)

func TestRequireAuthMiddleware(t *testing.T) {
	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Dummy handler that should not be called
		next := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			t.Error("next handler should not be called")
			return nil
		})

		if err := RequireAuthMiddleware(next)(w, req); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		user := db.User{ID: 1, Name: "test"}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		called := false
		next := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			called = true
			w.WriteHeader(http.StatusOK)
			return nil
		})

		if err := RequireAuthMiddleware(next)(w, req); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if !called {
			t.Error("next handler should have been called")
		}
	})
}
