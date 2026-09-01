package middleware

import (
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

var Buckets []*TokenBucket

func NewTokenBucket(ip string, capacity, refillRate int) *TokenBucket {
	var tb TokenBucket

	tb.ip = ip
	tb.capacity = capacity
	tb.tokens = float64(capacity)
	tb.refillRate = refillRate
	tb.lastRefillTime = time.Now()

	return &tb
}

func FindBucket(ip string, bk []*TokenBucket) *TokenBucket {
	for _, bucket := range bk {
		if bucket.ip == ip {
			return bucket
		}
	}
	return nil
}
func GetBucket(ip string) *TokenBucket {
	bucket := FindBucket(ip, Buckets)
	if bucket == nil {
		bucket = NewTokenBucket(ip, 5, 1)
		Buckets = append(Buckets, bucket)
	}

	return bucket
}
func AllowRequest(ip string) bool {
	bucket := GetBucket(ip)

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

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			http.Error(w, "Invalid client address", http.StatusInternalServerError)
			return
		}

		if !AllowRequest(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
