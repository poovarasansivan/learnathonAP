package roles

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func GetRole(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT id, role FROM m_users WHERE id = ?", requestData.UserID)
	var events UsersRole
	err := row.Scan(&events.ID, &events.UserRole)

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
		Events UsersRole `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func GetRoleC(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT id, role FROM m_users WHERE id = ?", requestData.UserID)
	var events UsersRoleC
	err := row.Scan(&events.ID, &events.UserRole)

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
		Events UsersRoleC `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCRole(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT name FROM m_users WHERE addcategory_role = '1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []UserCRole
	for rows.Next() {
		var user UserCRole
		if err := rows.Scan(&user.Name); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []UserCRole `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCategoryCountR(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var categories []CategoryCountR
	var temp CategoryCountR

	// Parse the request to get the 'id'
	var input InputR
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}

	row, err := config.Database.Query("SELECT COUNT(*) AS category_count FROM event_register WHERE event_category_id=? AND STATUS='1'", input.Id)

	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}

	for row.Next() {
		err := row.Scan(&temp.CRcount)
		if err != nil {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
			function.Response(w, response)
			return
		}

		tempRow := CategoryCountR{
			CRcount: temp.CRcount,
		}
		categories = append(categories, tempRow)
	}
	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}