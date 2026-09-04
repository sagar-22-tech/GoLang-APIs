package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"api/m/v2/apis/users/models"

	"github.com/redis/go-redis/v9"
)

func UserHandler(rdb *redis.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed",
			})
			return
		}

		search := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("search")))

		if search == "" {

			file, err := os.Open("apis/users/db.json")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Database file missing or unavailable",
				})
				return
			}
			defer file.Close()

			var users []models.User

			if err := json.NewDecoder(file).Decode(&users); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Malformed database structure",
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(users)
			return
		}

		cacheKey := "user:search:" + search

		cached, err := rdb.Get(r.Context(), cacheKey).Result()

		if err == nil {

			fmt.Println("CACHE HIT")

			var matchedUsers []models.User

			if err := json.Unmarshal([]byte(cached), &matchedUsers); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid cached data",
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(matchedUsers)
			return

		} else if err != redis.Nil {

			fmt.Println("Redis error:", err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Redis error",
			})
			return
		}

		fmt.Println("CACHE MISS")

		file, err := os.Open("apis/users/db.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Database file missing or unavailable",
			})
			return
		}
		defer file.Close()

		var users []models.User

		if err := json.NewDecoder(file).Decode(&users); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Malformed database structure",
			})
			return
		}

		var matchedUsers []models.User

		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Name.Firstname), search) {
				matchedUsers = append(matchedUsers, user)
			}
		}

		if len(matchedUsers) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "User not found",
			})
			return
		}

		data, err := json.Marshal(matchedUsers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to encode user data",
			})
			return
		}

		err = rdb.Set(
			r.Context(),
			cacheKey,
			data,
			0,
		).Err()

		if err != nil {
			fmt.Println("Redis SET error:", err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(matchedUsers)
	}
}

func UserHandlerID(rdb *redis.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method Not allowed",
			})
			return
		}

		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "Missing user ID", http.StatusBadRequest)
			return
		}

		targetID, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid ID format: must be an integer", http.StatusBadRequest)
			return
		}

		targetID = ((targetID - 1) % 10) + 1

		// Redis key for this user
		cacheKey := "user:id:" + strconv.Itoa(targetID)

		// Check Redis first
		cached, err := rdb.Get(r.Context(), cacheKey).Result()

		if err == nil {

			// CACHE HIT
			fmt.Println("CACHE HIT")

			var matchedUser models.User

			if err := json.Unmarshal([]byte(cached), &matchedUser); err != nil {
				http.Error(w, "Invalid cached user data", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(matchedUser)
			return

		} else if err != redis.Nil {

			// Actual Redis error
			fmt.Println("Redis error:", err)

			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}

		// CACHE MISS
		fmt.Println("CACHE MISS")

		file, err := os.OpenFile("apis/users/db.json", os.O_RDONLY, 0)

		if err != nil {
			http.Error(
				w,
				"Database file missing or unavailable",
				http.StatusInternalServerError,
			)
			return
		}

		defer file.Close()

		decoder := json.NewDecoder(file)

		if _, err := decoder.Token(); err != nil {
			http.Error(
				w,
				"Malformed database structure",
				http.StatusInternalServerError,
			)
			return
		}

		var matchedUser *models.User

		for decoder.More() {

			var u models.User

			if err := decoder.Decode(&u); err != nil {
				http.Error(
					w,
					"Error reading database entry",
					http.StatusInternalServerError,
				)
				return
			}

			if u.ID == targetID {
				matchedUser = &u
				break
			}
		}

		if matchedUser == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Convert user to JSON
		data, err := json.Marshal(matchedUser)

		if err != nil {
			http.Error(
				w,
				"Failed to encode user data",
				http.StatusInternalServerError,
			)
			return
		}

		// Store in Redis permanently
		err = rdb.Set(
			r.Context(),
			cacheKey,
			data,
			0,
		).Err()

		if err != nil {
			// Don't fail the API just because caching failed.
			fmt.Println("Redis SET error:", err)
		}

		// Return user
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(matchedUser)
	}
}
