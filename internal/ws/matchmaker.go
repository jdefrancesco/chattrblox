// internal/ws/matchmaker.go
package ws

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Matchmaker struct {
	Redis    *redis.Client
	QueueKey string
	Clients  map[uuid.UUID]*Client
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

	id1, _ := uuid.Parse(users[0])
	id2, _ := uuid.Parse(users[1])

	c1 := m.Clients[id1]
	c2 := m.Clients[id2]

	if c1 != nil && c2 != nil {
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
	}
}
