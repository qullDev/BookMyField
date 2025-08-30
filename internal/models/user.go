package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Password  string    `gorm:"not null" json:"-"`                 // hidden dari JSON
	Role      string    `gorm:"not null;default:user" json:"role"` // user or admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
