package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Represents a client.
type Client struct {
	ID     uuid.UUID
	Conn   *websocket.Conn
	SendCh chan Message
	Hub    *Hub
}

// Our message. Simple JSON entity containing message type along with data.
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
	To   string `json:"to,omitempty"` // Partner user ID
}

// Read incoming messages.
func (c *Client) ReadLoop() {
	defer c.Conn.Close()
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			break
		}
		// Handle incoming messages.
		switch msg.Type {
		case "signal_offer", "signal_answer", "ice_candidate":
			c.relayToPartner(msg)
		case "disconnect":
			c.handleDisconnect()
		}
	}
}

// Semd JSON msg.
func (c *Client) WriteLoop() {
	defer c.Conn.Close()
	for msg := range c.SendCh {
		c.Conn.WriteJSON(msg)
	}
}

// Relay message to paired partner.
func (c *Client) relayToPartner(msg Message) {
	partnerID, err := uuid.Parse(msg.To)
	if err != nil {
		return
	}

	partner := c.Hub.Matchmaker.Clients[partnerID]
	if partner != nil {
		partner.SendCh <- Message{
			Type: msg.Type,
			Data: msg.Data,
			To:   c.ID.String(), // sender ID
		}
	}
}

// Handle disconnect.
func (c *Client) handleDisconnect() {
	// Remove user from matchmaking pool and notify partner
	delete(c.Hub.Matchmaker.Clients, c.ID)
	c.Conn.Close()
}
