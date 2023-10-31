package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Sample(w http.ResponseWriter, r *http.Request) {
	fmt.Print("buibib")
	response := map[string]interface{}{
		"success": false,
		"error":   "Invalid Email ID",
	}

	var err error
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to encode response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
