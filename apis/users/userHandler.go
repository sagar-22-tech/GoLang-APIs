package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"api/m/v2/apis/users/models"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}
	fileBytes, err := os.ReadFile("apis/users/db.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	var users []models.User
	err = json.Unmarshal(fileBytes, &users)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}
func UserHandler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method Not allowed",
		})
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

	file, err := os.OpenFile("apis/users/db.json", os.O_RDONLY, 0)
	if err != nil {
		http.Error(w, "Database file missing or unavailable", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		http.Error(w, "Malformed database structure", http.StatusInternalServerError)
		return
	}

	var matchedUser *models.User

	for decoder.More() {
		var u models.User
		if err := decoder.Decode(&u); err != nil {
			http.Error(w, "Error reading database entry", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(matchedUser)
}
