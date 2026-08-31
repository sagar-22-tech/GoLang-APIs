package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"example.com/m/v2/apis/greet"
	"example.com/m/v2/apis/health"
	"example.com/m/v2/apis/users"
	"github.com/rs/cors"
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", health.HealthHandler)
	mux.HandleFunc("/user", users.UserHandler)

	mux.HandleFunc("/greet", greet.GreetHandlerWithTime)

	c := cors.New(cors.Options{

		AllowedOrigins:   []string{"https://go-playground-weld.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	bucket := NewTokenBucket(5, 1)
	rateLimitedMux := RateLimitMiddleware(bucket, mux)

	handler := c.Handler(rateLimitedMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on Port %s\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		panic(err)
	}
}
