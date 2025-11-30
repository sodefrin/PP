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

const UserContextKey contextKey = "user"

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
					slog.Error("GetSession error", "error", err)
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
				slog.Error("GetUser error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
