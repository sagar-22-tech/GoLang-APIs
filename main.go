package main

import (
	"fmt"
	"net/http"
	"os"

	"api/m/v2/apis/greet"
	"api/m/v2/apis/health"
	"api/m/v2/apis/users"
	"api/m/v2/middleware"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	//Health api endpoint
	mux.HandleFunc("/health", health.HealthHandler)

	//Users api endpoint
	mux.HandleFunc("/users", users.UserHandler)
	mux.HandleFunc("/users/{id}", users.UserHandler2)

	mux.HandleFunc("/greet", greet.GreetHandlerWithTime)

	c := cors.New(cors.Options{

		AllowedOrigins:   []string{"https://go-playground-weld.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	bucket := middleware.NewTokenBucket(5, 1)
	rateLimitedMux := middleware.RateLimitMiddleware(bucket, mux)

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
