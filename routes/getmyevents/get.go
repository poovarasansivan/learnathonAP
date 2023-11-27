package getmyevents

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"log"
	"net/http"
)

func GetMyEvents(w http.ResponseWriter, r *http.Request) {
	// Parse request 	body
	var requestData struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow(
		"SELECT er.team_name, er.user_1, mu1.name AS user_1_name, er.user_2,mu2.name AS user_2_name, er.user_3,mu3.name AS user_3_name, muu.name AS event_incharge,muc.name AS cincharge, me.event_name, me.description AS edescription,er.event_category_id, mc.category_name AS cname, mc.description AS cdescription, me.event_date FROM `event_register` er INNER JOIN `event_categories` ec ON ec.`id` = er.`event_category_id` INNER JOIN m_category mc ON mc.id=er.event_category_id INNER JOIN m_events me ON me.id = ec.`event_id` INNER JOIN m_users muu ON muu.id=me.incharge LEFT JOIN `m_users` mu1 ON mu1.id = er.`user_1` LEFT JOIN m_users mu2 ON mu2.id = er.`user_2` LEFT JOIN m_users mu3 ON mu3.id = er.`user_3` INNER JOIN m_users muc ON muc.id=mc.incharge WHERE (er.user_1 =? OR er.user_2 =? OR er.user_3 =?)",
		requestData.UserID, requestData.UserID, requestData.UserID)

	var events MyEvents
	err := row.Scan(
		&events.TeamName, &events.User1, &events.User1_Name, &events.User2, &events.User2_Name, &events.User3,
		&events.User3_Name, &events.EIncharge, &events.CIncharge, &events.EventName, &events.Edesciption, &events.Category_id, &events.CategoryName,
		&events.CDescription, &events.EventDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Events not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Events MyEvents `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMyCategorys(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		User_ID string `json:"user_1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT event_category_id FROM event_register WHERE user_1=? AND STATUS='1'", requestData.User_ID)

	var categoryID int
	err := row.Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No category found for the user", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Event GetMyCategory `json:"event"`
	}{Event: GetMyCategory{CategoryID: categoryID}}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}