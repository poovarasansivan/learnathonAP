package categorydata

import (
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func InsertcategoryData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Category_Name string `json:"category_name"`
		Description   string `json:"description"`
		Max_Team      int    `json:"max_team"`
		Incharge      string `json:"incharge"`
		Created_by    string `json:"created__by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := config.Database.Exec("INSERT INTO m_category (category_name, description, max_team, incharge,status,created_by, created_at, updated_on) VALUES (?, ?, ?, ?, '1', ?, NOW(), NOW())",
		req.Category_Name, req.Description, req.Max_Team, req.Incharge, req.Created_by)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}


func GetCategoryName(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT id,category_name FROM m_category WHERE STATUS = '1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []QuestionCategory
	for rows.Next() {
		var user QuestionCategory
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []QuestionCategory `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}