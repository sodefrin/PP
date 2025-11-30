package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/sodefrin/PP/server/api"
	"github.com/sodefrin/PP/server/db"

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

	slog.Info("Database initialized (in-memory)")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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
