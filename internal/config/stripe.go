package config

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v76"
)

func InitStripe() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		log.Println("⚠️ STRIPE_SECRET_KEY not found, Stripe features will be disabled")
		return
	}

	stripe.Key = key
	log.Println("✅ Stripe initialized successfully")
}
