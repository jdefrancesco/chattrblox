package ws

import (
	"chatrblox/internal/middleware"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	Matchmaker *Matchmaker
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		ID:     userID,
		Conn:   conn,
		Hub:    h,
		SendCh: make(chan Message),
	}

	go client.ReadLoop()
	go client.WriteLoop()

	h.Matchmaker.Enqueue(client)
}
