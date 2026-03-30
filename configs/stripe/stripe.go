package configs

import (
	"log/slog"
	"os"

	"github.com/stripe/stripe-go/v79"
)

func InitStripe() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		slog.Warn("STRIPE_SECRET_KEY not set — payment processing disabled")
		return
	}

	stripe.Key = key
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: stripe.LevelWarn, // suppress verbose SDK logs in production
	}

	slog.Info("Stripe initialized", "mode", stripeMode(key))
}

func StripeEnabled() bool {
	return stripe.Key != ""
}

func stripeMode(key string) string {
	if len(key) > 7 && key[:7] == "sk_live" {
		return "live"
	}
	return "test"
}
