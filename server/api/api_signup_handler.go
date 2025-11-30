package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sodefrin/PP/server/api/dto"
	"github.com/sodefrin/PP/server/db"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dto.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Password == "" {
			http.Error(w, "Name and password are required", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Bcrypt error", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		params := db.CreateUserParams{
			Name:         req.Name,
			PasswordHash: string(hash),
		}

		user, err := queries.CreateUser(context.Background(), params)
		if err != nil {
			slog.Error("CreateUser error", "error", err)
			http.Error(w, "Failed to create user (name might be taken)", http.StatusConflict)
			return
		}

		// Create session
		sessionID := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		sessionParams := db.CreateSessionParams{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
		}

		_, err = queries.CreateSession(context.Background(), sessionParams)
		if err != nil {
			slog.Error("CreateSession error", "error", err)
			// Don't fail the request, just log error. User is created.
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Expires:  expiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(userJSON); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	}
}
