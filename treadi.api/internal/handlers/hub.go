package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"treadi.api/internal/hub"
)

type HubHandler struct {
	hub      *hub.Hub
	deadline time.Time
	logger   *slog.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Event struct {
	Event   string `json:"event"`
	Payload *json.RawMessage
}

func (hh HubHandler) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			hh.logger.Debug("failed to handshake")
		}
		return
	}
	client := hub.NewClient(hh.hub, conn)
	client.Register()
	go client.Write()
	go client.Read()
}

func NewHubHandler(hub *hub.Hub, logger *slog.Logger) HubHandler {
	child := logger.With(
		slog.Group("context",
			slog.String("function", "hub_handler"),
		),
	)
	return HubHandler{
		hub:    hub,
		logger: child,
	}
}
