package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"api/m/v2/apis/greet"
	"api/m/v2/apis/health"
	"api/m/v2/apis/users"
	"api/m/v2/middleware"

	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()

	// Health API endpoint
	mux.HandleFunc("/health", health.HealthHandler)

	// Users API endpoint
	mux.HandleFunc("/users", users.UserHandler(rdb))
	mux.HandleFunc("/users/{id}", users.UserHandlerID(rdb))

	// Greet API endpoint
	mux.HandleFunc("/greet", greet.GreetHandlerWithTime)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://go-playground-weld.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Create Leaky Bucket
	bucket := middleware.NewLeakyBucket(
		5,             // capacity
		1*time.Second, // leak rate
		5*time.Second, // request expiry
	)

	// Start background leak worker
	go bucket.Leak()

	// Rate limit middleware
	rateLimitedHandler := middleware.RateLimitMiddleware(
		bucket,
		mux,
	)

	// CORS wraps rate limiter
	handler := c.Handler(rateLimitedHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on Port %s\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		panic(err)
	}
}
