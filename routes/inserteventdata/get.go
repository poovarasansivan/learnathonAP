package inserteventdata

import (
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"net/http"
	"time"
)

func InsertEventData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EventName   string `json:"event_name"`
		Description string `json:"description"`
		Event_date  string `json:"event_date"`
		Incharge    string `json:"incharge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	eventDate, err := time.Parse("2006-01-02", req.Event_date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = config.Database.Exec("INSERT INTO m_events (event_name, description, event_date, incharge, status, created_at, updated_on) VALUES (?, ?, ?, ?, '1', NOW(), NOW())",
		req.EventName, req.Description, eventDate, req.Incharge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}