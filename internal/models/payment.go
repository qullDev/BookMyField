package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null" json:"booking_id"`
	Booking   Booking   `gorm:"foreignKey:BookingID"`

	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`        // pending, paid, failed
	StripeRefID string    `json:"stripe_ref_id"` // session ID atau payment intent ID
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
