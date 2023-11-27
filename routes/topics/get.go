package topics

import (
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"net/http"
	
)

func GetTopics(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			response := map[string]interface{}{
				"success": false,
				"error":   "Internal Server Error",
			}
			function.Response(w, response)
		}
	}()

	var response map[string]interface{}
	var categories []Topics
	var input CategoryIdInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	rows, err := config.Database.Query("SELECT topics FROM m_topics WHERE STATUS='1' AND category_id=?", input.CId)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category Topics
		err := rows.Scan(&category.TopicsName)
		if err != nil {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
			function.Response(w, response)
			return
		}
		categories = append(categories, category)
	}

	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}