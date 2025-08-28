package config

import (
	"os"

	"github.com/stripe/stripe-go/v72"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}
