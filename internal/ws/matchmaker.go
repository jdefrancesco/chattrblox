package ws

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Matchmaker struct has necesarry data to track active users.
// It relies om Redis for a fast datastore.
type Matchmaker struct {
	Redis    *redis.Client
	QueueKey string
	Clients  map[uuid.UUID]*Client
	Sessions map[uuid.UUID]uuid.UUID
	DB       *gorm.DB
}

func NewMatchmaker(redisClient *redis.Client) *Matchmaker {
	return &Matchmaker{
		Redis:    redisClient,
		QueueKey: "matchmaking_queue",
		Clients:  make(map[uuid.UUID]*Client),
	}
}

func (m *Matchmaker) Enqueue(client *Client) {
	m.Clients[client.ID] = client
	m.Redis.LPush(context.Background(), m.QueueKey, client.ID.String())

	m.TryMatch()
}

func (m *Matchmaker) TryMatch() {
	ctx := context.Background()

	users, err := m.Redis.RPopCount(ctx, m.QueueKey, 2).Result()
	if err != nil || len(users) < 2 {
		return
	}

	id1, err1 := uuid.Parse(users[0])
	id2, err2 := uuid.Parse(users[1])
	if err1 != nil || err2 != nil {
		log.Println("Failed to parse user IDs")
		return
	}

	c1 := m.Clients[id1]
	c2 := m.Clients[id2]

	if c1 != nil && c2 != nil {
		// Store session relationship
		m.Sessions[id1] = id2
		m.Sessions[id2] = id1

		log.Printf("Matched %s with %s", id1.String(), id2.String())

		notify := func(c *Client, partner uuid.UUID) {
			msg := Message{
				Type: "match_found",
				Data: partner.String(),
			}
			c.SendCh <- msg
		}
		notify(c1, c2.ID)
		notify(c2, c1.ID)

	} else {
		log.Println("One or both users disconnected before matching")

		// Requeue connected client if oone exists.
		if c1 != nil {
			m.Enqueue(c1)
		}
		if c2 != nil {
			m.Enqueue(c2)
		}
	}
}
