package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity       int
	tokens         float64
	refillRate     int
	lastRefillTime time.Time
	mu             sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	var tb TokenBucket
	tb.capacity = capacity
	tb.tokens = float64(capacity)
	tb.refillRate = refillRate
	tb.lastRefillTime = time.Now()

	return &tb

}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	elapsedTime := now.Sub(tb.lastRefillTime)
	tokensToAdd := elapsedTime.Seconds() * float64(tb.refillRate)
	tb.lastRefillTime = now
	tb.tokens += tokensToAdd

	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(bucket *TokenBucket, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !bucket.Allow() {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{
				"message": "Retry after sometime",
				"error":   "Too may requests !",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
