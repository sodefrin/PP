package lib

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/sodefrin/PP/server/db"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				// No session cookie, continue without user
				next.ServeHTTP(w, r)
				return
			}

			sessionID := cookie.Value
			session, err := queries.GetSession(r.Context(), sessionID)
			if err != nil {
				if err != sql.ErrNoRows {
					slog.ErrorContext(r.Context(), "GetSession error", "error", err)
				}
				// Invalid session, continue without user
				next.ServeHTTP(w, r)
				return
			}

			if time.Now().After(session.ExpiresAt) {
				// Session expired
				next.ServeHTTP(w, r)
				return
			}

			user, err := queries.GetUser(r.Context(), session.UserID)
			if err != nil {
				slog.ErrorContext(r.Context(), "GetUser error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, ok := r.Context().Value(userContextKey).(db.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}
		return next(w, r)
	}
}
