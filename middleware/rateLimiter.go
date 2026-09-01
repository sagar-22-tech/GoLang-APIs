package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	ip             string
	capacity       int
	tokens         float64
	refillRate     int
	lastRefillTime time.Time
	mu             sync.Mutex
}

type RateLimiter struct {
	buckets    map[string]*TokenBucket
	capacity   int
	refillRate int
	mu         sync.Mutex
}

func NewRateLimiter(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func NewTokenBucket(capacity, refillRate int, ip string) *TokenBucket {
	return &TokenBucket{
		ip:             ip,
		capacity:       capacity,
		tokens:         float64(capacity),
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (rl *RateLimiter) Allow(ip string) bool {

	rl.mu.Lock()

	bucket, exists := rl.buckets[ip]

	if !exists {
		bucket = NewTokenBucket(rl.capacity, rl.refillRate, ip)
		rl.buckets[ip] = bucket
	}

	rl.mu.Unlock()

	return bucket.Allow()
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

func RateLimitMiddleware(
	limiter *RateLimiter,
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			http.Error(
				w,
				"Invalid client address",
				http.StatusInternalServerError,
			)
			return
		}

		if !limiter.Allow(ip) {
			w.Header().Set("Content-Type", "application/json")

			response := map[string]string{
				"message": "Retry after sometime",
				"error":   "Too many requests!",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
