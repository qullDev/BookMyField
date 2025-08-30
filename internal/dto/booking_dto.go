package dto

import "time"

// CreateBookingRequest represents the request body for creating a booking
type CreateBookingRequest struct {
	FieldID   string    `json:"field_id" binding:"required" example:"c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d"`
	StartTime time.Time `json:"start_time" binding:"required" example:"2024-09-15T10:00:00Z"`
	EndTime   time.Time `json:"end_time" binding:"required" example:"2024-09-15T12:00:00Z"`
}

// CreateFieldRequest represents the request body for creating a field
type CreateFieldRequest struct {
	Name     string  `json:"name" binding:"required" example:"Lapangan Futsal A"`
	Location string  `json:"location" binding:"required" example:"Jakarta"`
	Price    float64 `json:"price" binding:"required" example:"200000"`
}

// UpdateFieldRequest represents the request body for updating a field
type UpdateFieldRequest struct {
	Name     string  `json:"name" binding:"required" example:"Lapangan Futsal A Updated"`
	Location string  `json:"location" binding:"required" example:"Jakarta Barat"`
	Price    float64 `json:"price" binding:"required" example:"250000"`
}

// CreateCheckoutSessionRequest represents the request body for creating a Stripe checkout session
type CreateCheckoutSessionRequest struct {
	BookingID string `json:"booking_id" binding:"required" example:"c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d"`
}

// CreateCheckoutSessionResponse represents the response for creating a Stripe checkout session
type CreateCheckoutSessionResponse struct {
	SessionID  string `json:"session_id" example:"cs_test_..."`
	SessionURL string `json:"session_url" example:"https://checkout.stripe.com/pay/cs_test_..."`
}

// CancelBookingResponse represents the response for cancelling a booking
type CancelBookingResponse struct {
	Message      string `json:"message" example:"Booking cancelled and payment refunded successfully"`
	RefundID     string `json:"refund_id,omitempty" example:"re_..."`
	RefundStatus string `json:"refund_status,omitempty" example:"succeeded"`
}
