package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Report struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReporterID uuid.UUID
	ReportedID uuid.UUID
	Reason     string
	CreatedAt  time.Time
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	r.CreatedAt = time.Now()
	return
}
