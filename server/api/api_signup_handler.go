package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sodefrin/PP/server/api/dto"
	"github.com/sodefrin/PP/server/db"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := Queries.CreateUser(context.Background(), params)
	if err != nil {
		slog.Error("CreateUser error", "error", err)
		http.Error(w, "Failed to create user (name might be taken)", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
