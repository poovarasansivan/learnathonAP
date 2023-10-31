package function

import (
	"encoding/json"
	"log"
	"net/http"
)

func Response(w http.ResponseWriter, response map[string]interface{}) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to encode response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
