package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	User1ID        uuid.UUID  `gorm:"type:uuid;not null"`
	User2ID        uuid.UUID  `gorm:"type:uuid;not null"`
	StartedAt      time.Time  `gorm:"not null"`
	EndedAt        time.Time  `gorm:"not null"`
	DisconnectedBy *uuid.UUID `gorm:"type:uuid"` // nullable
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	s.StartedAt = time.Now()
	return
}
