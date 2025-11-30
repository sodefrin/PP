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
	"github.com/sodefrin/PP/server/lib"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(queries *db.Queries) lib.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return nil
		}

		var req dto.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return nil
		}

		if req.Name == "" || req.Password == "" {
			http.Error(w, "Name and password are required", http.StatusBadRequest)
			return nil
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		params := db.CreateUserParams{
			Name:         req.Name,
			PasswordHash: string(hash),
		}

		user, err := queries.CreateUser(context.Background(), params)
		if err != nil {
			http.Error(w, "Failed to create user (name might be taken)", http.StatusConflict)
			return nil
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
			slog.ErrorContext(r.Context(), "CreateSession error", "error", err)
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

		resp := dto.User{
			ID:   user.ID,
			Name: user.Name,
		}

		userJSON, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(userJSON); err != nil {
			return err
		}
		return nil
	}
}
