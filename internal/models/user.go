package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"primary_key"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"not null;unique"`
	Password   string    `gorm:"not null"`
	Role       string    `gorm:"not null;default:user"` // user or admin
	CreadtedAt time.Time
}

func (u User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
