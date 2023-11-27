package categorydetails

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)

func GetDetail(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var categories CategoryDetails
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	err = config.Database.QueryRow("SELECT mc.id,mc.category_name,mc.description,mu.name,mc.max_team,mc.registration FROM event_categories ec INNER JOIN m_category mc ON mc.id = ec.category_id INNER JOIN m_events mee ON mee.id = ec.event_id INNER JOIN m_users mu ON mu.id = mc.incharge WHERE ec.status = '1' AND ec.category_id = ?", input.Id).Scan(&categories.Id, &categories.Name, &categories.Description, &categories.InchargeName, &categories.MaxTeam, &categories.Registerstatus)

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No Request",
			}
		} else {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		function.Response(w, response)
		return
	}

	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}