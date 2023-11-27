package registerdata

import (
	"encoding/json"
	"fmt"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func InsertData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TeamName        string `json:"teamName"`
		EventCategoryID int    `json:"eventCategoryID"`
		User1           string `json:"user1"`
		User2           string `json:"user2"`
		User3           string `json:"user3"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := config.Database.Exec("INSERT INTO event_register (team_name, event_category_id, user_1, user_2, user_3, status, created_by, created_on, updated_on) VALUES (?, ?, ?, ?, ?, '1', ?, NOW(), NOW())",
		req.TeamName, req.EventCategoryID, req.User1, req.User2, req.User3, req.User1)
	fmt.Print(req.User1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}


func CheckTeam(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		RollNo string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var isRegistered bool
	rows, err := config.Database.Query(`
        SELECT er.user_1, er.user_2, er.user_3 
        FROM event_register er 
        WHERE er.status='1' 
            AND (er.user_1=? OR er.user_2=? OR er.user_3=?)`,
		requestData.RollNo, requestData.RollNo, requestData.RollNo)

	if err != nil {

		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var team CTeam
		err := rows.Scan(&team.User1, &team.User2, &team.User3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		isRegistered = true
		break
	}
	response := struct {
		IsRegistered bool `json:"isRegistered"`
	}{IsRegistered: isRegistered}
	fmt.Print(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

