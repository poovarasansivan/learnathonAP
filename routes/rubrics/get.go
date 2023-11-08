package rubrics

import (
	"database/sql"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)

func GetRubrics(w http.ResponseWriter, r *http.Request) {
	var data []Criteria
	var dataRe []Rubrics

	var temp Criteria
	var tempRub Rubrics
	var response map[string]interface{}

	criteria, err := config.Database.Query("SELECT id,NAME FROM m_rubrics_criteria WHERE STATUS ='1'")
	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No Criteria Found",
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

	for criteria.Next() {
		criteria.Scan(&temp.CriteriaID, &temp.CriteriaName)
		rubrics, _ := config.Database.Query("SELECT id,NAME FROM m_rubrics_questions WHERE criteria_id =? AND  STATUS ='1'", temp.CriteriaID)

		for rubrics.Next() {
			rubrics.Scan(&tempRub.RubricsID, &tempRub.RubricsName)
			dataRe = append(dataRe, tempRub)
		}
		temp.Rubrics = dataRe
		data = append(data, temp)
	}
	response = map[string]interface{}{
		"success": true,
		"data":    data,
	}
	function.Response(w, response)
}
