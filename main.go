package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"puyo-server/server/db"

	"github.com/gorilla/websocket"
	_ "modernc.org/sqlite"
)

//go:embed public/*
var content embed.FS

//go:embed server/db/schema.sql
var schema string

var queries *db.Queries
var dbConn *sql.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		log.Printf("Received: %s", p)

		// Echo message back
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
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

	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/ws", wsHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
