package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct holds information of registered users.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	IsAdmin      bool      `gorm:"default:false"`
	CreatedAt    time.Time
	LastLoginAt  *time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	return
}
