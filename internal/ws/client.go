package ws

import (
	"chatrblox/internal/models"
	"log"

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
		case "report":
			c.handleReport(msg)
		}
	}
}

// Send JSON msg.
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

// Gracefully handle disconnection.
func (c *Client) handleDisconnect() {
	partnerID := c.Hub.Matchmaker.FindPartnerID(c.ID)

	if partnerID != uuid.Nil {
		partner := c.Hub.Matchmaker.Clients[partnerID]
		if partner != nil {
			partner.SendCh <- Message{
				Type: "disconnect",
				Data: "partner left",
			}
			// Re-queue the partner
			c.Hub.Matchmaker.Enqueue(partner)
		}
	}

	// Cleanup session entry.
	delete(c.Hub.Matchmaker.Sessions, c.ID)
	delete(c.Hub.Matchmaker.Sessions, partnerID)
}

// Handle report button.
func (c *Client) handleReport(msg Message) {
	reportedID, err := uuid.Parse(msg.To)
	if err != nil {
		return
	}

	report := models.Report{
		ReporterID: c.ID,
		ReportedID: reportedID,
		Reason:     msg.Data,
	}

	if err := c.Hub.Matchmaker.DB.Create(&report).Error; err != nil {
		log.Println("Failed to log report:", err)
	}
}
