package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("Upgrade error", "error", err)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				slog.Error("Failed to close websocket connection", "error", err)
			}
		}()

		slog.Info("Client connected")

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Read error", "error", err)
				return
			}
			slog.Info("Received message", "payload", string(p))

			// Echo message back
			if err := conn.WriteMessage(messageType, p); err != nil {
				slog.Error("Write error", "error", err)
				return
			}
		}
	}
}
