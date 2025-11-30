package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sodefrin/PP/server/db"
	"github.com/sodefrin/PP/server/lib"
)

func TestMeHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		user := db.User{
			ID:   1,
			Name: "testuser",
		}

		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		// Inject user into context (simulating AuthMiddleware)
		ctx := context.WithValue(req.Context(), lib.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		MeHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.ID != user.ID {
			t.Errorf("expected ID %d, got %d", user.ID, resp.ID)
		}
		if resp.Name != user.Name {
			t.Errorf("expected Name %s, got %s", user.Name, resp.Name)
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/me", nil)
		w := httptest.NewRecorder()
		MeHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}

func TestRequireAuthMiddleware(t *testing.T) {
	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Dummy handler that should not be called
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("next handler should not be called")
		})

		lib.RequireAuthMiddleware(next)(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		user := db.User{ID: 1, Name: "test"}
		ctx := context.WithValue(req.Context(), lib.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		lib.RequireAuthMiddleware(next)(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if !called {
			t.Error("next handler should have been called")
		}
	})
}
