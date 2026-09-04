package greet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GreetHandlerWithTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method Not Allowed",
		})

		return
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	currentTimeIST := time.Now().In(loc)

	response := map[string]string{
		"message": "Hello,World!",
		"time":    currentTimeIST.String(),
		"method":  r.Method,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
