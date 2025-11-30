package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sodefrin/PP/server/lib"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler() lib.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
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
				return err
			}
			slog.Info("Received message", "payload", string(p))

			// Echo message back
			if err := conn.WriteMessage(messageType, p); err != nil {
				return err
			}
		}
	}
}
