package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/sodefrin/PP/server/api"
	"github.com/sodefrin/PP/server/db"
	"github.com/sodefrin/PP/server/lib"

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
	// Use otelsql to open database with tracing
	dbConn, err = otelsql.Open("sqlite", "file::memory:?cache=shared",
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithSQLCommenter(true),
	)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	// Register DB stats metrics (optional but good practice)
	if err := otelsql.RegisterDBStatsMetrics(dbConn, otelsql.WithAttributes(semconv.DBSystemSqlite)); err != nil {
		slog.Error("Failed to register DB stats metrics", "error", err)
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

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(publicFS)))

	mux.HandleFunc("/api/health", api.HealthHandler)
	mux.HandleFunc("/ws", api.WsHandler)
	mux.HandleFunc("/api/signup", api.SignupHandler)
	mux.HandleFunc("/api/signin", api.SigninHandler)

	// Wrap mux with logging middleware
	handler := lib.LoggingMiddleware(mux)

	// Wrap with OpenTelemetry
	handler = otelhttp.NewHandler(handler, "server")

	// Initialize Tracer
	tp, err := lib.InitTracer()
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracer", "error", err)
		}
	}()

	port := ":8080"
	slog.Info("Server starting", "port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		slog.Error("ListenAndServe failed", "error", err)
		os.Exit(1)
	}
}
