package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method Not Allowed",
			})
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Go API is operational",
		})
	})

	mux.HandleFunc("/api/v1/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method Not Allowed",
			})

			return
		}

		var requestData map[string]any

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid JSON body",
			})

			return
		}

		name, ok := requestData["name"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "name must be a string",
			})

			return
		}

		response := map[string]string{
			"message": "Hello, " + name + "!",
			"method":  r.Method,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	c := cors.New(cors.Options{

		AllowedOrigins:   []string{"https://go-playground-weld.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on Port %s\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		panic(err)
	}
}
