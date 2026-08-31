package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	ID   int `json:"id"`
	Name struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	} `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name   string `json:"name"`
		Post   string `json:"post"`
		Salary string `json:"salary"`
	} `json:"company"`
}

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
	var users []User
	err = json.Unmarshal(fileBytes, &users)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}
