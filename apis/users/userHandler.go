package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
