package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

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
		log.Fatal(err)
	}

	// Execute schema
	if _, err := dbConn.Exec(schema); err != nil {
		log.Fatal(err)
	}

	queries = db.New(dbConn)

	// Test query
	ctx := context.Background()
	stat, err := queries.CreateGameStat(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database initialized (in-memory). Created stat ID: %d", stat.ID)
}

func main() {
	initDB()
	defer dbConn.Close()

	// Serve static files from embedded filesystem
	publicFS, err := fs.Sub(content, "public")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(publicFS)))

	http.HandleFunc("/api/health", api.HealthHandler)
	http.HandleFunc("/ws", api.WsHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
