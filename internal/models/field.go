package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Field struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Location  string    `gorm:"type:varchar(255);not null" json:"location"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// generate UUID otomatis sebelum create
func (f *Field) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
