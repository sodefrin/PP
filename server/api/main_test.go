package api

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"puyo-server/server/db"

	_ "modernc.org/sqlite"
)

func TestMain(m *testing.M) {
	// Setup DB
	dbConn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	// Read schema file
	schemaBytes, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		slog.Error("Failed to read schema file", "error", err)
		os.Exit(1)
	}

	// Execute schema
	if _, err := dbConn.Exec(string(schemaBytes)); err != nil {
		slog.Error("Failed to execute schema", "error", err)
		os.Exit(1)
	}

	// Initialize queries
	Queries = db.New(dbConn)

	code := m.Run()

	dbConn.Close()
	os.Exit(code)
}
