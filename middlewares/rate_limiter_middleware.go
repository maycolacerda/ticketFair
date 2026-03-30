// middlewares/rate_limiter.go
package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	tokens    float64
	lastRefil time.Time
	mu        sync.Mutex
}

func (b *bucket) allow(rate float64, capacity float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefil).Seconds()
	b.lastRefil = now

	// Refil tokens based on elapsed time
	b.tokens += elapsed * rate
	if b.tokens > capacity {
		b.tokens = capacity
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

type rateLimiterStore struct {
	buckets map[string]*bucket
	mu      sync.RWMutex
}

func newStore() *rateLimiterStore {
	s := &rateLimiterStore{
		buckets: make(map[string]*bucket),
	}
	// Cleanup stale buckets every 5 minutes
	go s.cleanup()
	return s
}

func (s *rateLimiterStore) get(ip string) *bucket {
	s.mu.RLock()
	b, ok := s.buckets[ip]
	s.mu.RUnlock()

	if ok {
		return b
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock
	if b, ok = s.buckets[ip]; ok {
		return b
	}

	b = &bucket{
		tokens:    0,
		lastRefil: time.Now(),
	}
	s.buckets[ip] = b
	return b
}

func (s *rateLimiterStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		for ip, b := range s.buckets {
			b.mu.Lock()
			// Remove buckets inactive for more than 10 minutes
			if time.Since(b.lastRefil) > 10*time.Minute {
				delete(s.buckets, ip)
			}
			b.mu.Unlock()
		}
		s.mu.Unlock()
	}
}

type RateLimitConfig struct {
	Rate     float64 // tokens per second
	Capacity float64 // burst size
	Message  string
}

var (
	// Auth endpoints — strict: 5 requests, refills 1 per 3 seconds
	AuthRateLimit = RateLimitConfig{
		Rate:     0.33,
		Capacity: 5,
		Message:  "too many login attempts — try again later",
	}

	// Registration — moderate: 3 requests, refills 1 per 10 seconds
	RegisterRateLimit = RateLimitConfig{
		Rate:     0.1,
		Capacity: 3,
		Message:  "too many registration attempts — try again later",
	}

	// Public endpoints — relaxed: 30 requests, refills 10 per second
	PublicRateLimit = RateLimitConfig{
		Rate:     10,
		Capacity: 30,
		Message:  "too many requests — slow down",
	}

	// Admin endpoints — moderate: 20 requests, refills 5 per second
	AdminRateLimit = RateLimitConfig{
		Rate:     5,
		Capacity: 20,
		Message:  "too many requests",
	}

	// Verification codes — strict: 3 requests, refills 1 per minute
	VerifyRateLimit = RateLimitConfig{
		Rate:     0.016,
		Capacity: 3,
		Message:  "too many verification attempts — try again later",
	}
)

// stores are shared per config type so IPs are tracked across requests
var stores = map[string]*rateLimiterStore{}
var storesMu sync.Mutex

func getStore(key string) *rateLimiterStore {
	storesMu.Lock()
	defer storesMu.Unlock()

	if s, ok := stores[key]; ok {
		return s
	}
	s := newStore()
	stores[key] = s
	return s
}

func RateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	store := getStore(cfg.Message) // use message as unique key per config

	return func(c *gin.Context) {
		ip := c.ClientIP()
		b := store.get(ip)

		if !b.allow(cfg.Rate, cfg.Capacity) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": cfg.Message,
			})
			return
		}

		c.Next()
	}
}
