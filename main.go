package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"puyo-server/server/api"
	"puyo-server/server/db"

	_ "modernc.org/sqlite"
)

//go:embed public/*
var content embed.FS

//go:embed server/db/schema.sql
var schema string

var queries *db.Queries
var dbConn *sql.DB

func initDB() {
	var err error
	dbConn, err = sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	// Execute schema
	if _, err := dbConn.Exec(schema); err != nil {
		slog.Error("Failed to execute schema", "error", err)
		os.Exit(1)
	}

	queries = db.New(dbConn)
	api.Queries = queries

	// Test query
	ctx := context.Background()
	stat, err := queries.CreateGameStat(ctx)
	if err != nil {
		slog.Error("Failed to create game stat", "error", err)
		os.Exit(1)
	}
	slog.Info("Database initialized (in-memory)", "stat_id", stat.ID)
}

func main() {
	initDB()
	defer dbConn.Close()

	// Serve static files from embedded filesystem
	publicFS, err := fs.Sub(content, "public")
	if err != nil {
		slog.Error("Failed to get public FS", "error", err)
		os.Exit(1)
	}
	http.Handle("/", http.FileServer(http.FS(publicFS)))

	http.HandleFunc("/api/health", api.HealthHandler)
	http.HandleFunc("/ws", api.WsHandler)
	http.HandleFunc("/api/signup", api.SignupHandler)
	http.HandleFunc("/api/signin", api.SigninHandler)

	port := ":8080"
	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error("ListenAndServe failed", "error", err)
		os.Exit(1)
	}
}
