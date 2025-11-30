package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sodefrin/PP/server/api/dto"
	"github.com/sodefrin/PP/server/db"
	"golang.org/x/crypto/bcrypt"
)

func SigninHandler(queries *db.Queries) http.HandlerFunc {
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

		user, err := queries.GetUserByName(context.Background(), req.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			slog.Error("GetUserByName error", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		userJSON, err := json.Marshal(user)
		if err != nil {
			slog.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(userJSON); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	}
}
