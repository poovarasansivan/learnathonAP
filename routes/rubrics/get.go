package rubrics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)

func GetRubrics(w http.ResponseWriter, r *http.Request) {
	var data []Criteria

	var temp Criteria
	var tempRub Rubrics
	var response map[string]interface{}

	criteria, err := config.Database.Query("SELECT id,NAME FROM m_rubrics_criteria WHERE STATUS ='1' ")
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

		var dataRe []Rubrics

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

func InsertRubricsData(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		Question_id int    `json:"question_id"`
		Criteria_ID int    `json:"criteria_id"`
		Rubrics_ID  int    `json:"selected"`
		Created_by  string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.Database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, question := range req {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM rubrics_log WHERE question_id=? AND criteria_id=? AND created_by=?",
			question.Question_id, question.Criteria_ID, question.Created_by).Scan(&count)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			// Data exists, perform an update
			_, err = tx.Exec("UPDATE rubrics_log SET rubrics_id=? WHERE question_id=? AND criteria_id=? AND created_by=?",
				question.Rubrics_ID, question.Question_id, question.Criteria_ID, question.Created_by)
		} else {
			// Data does not exist, perform an insert
			_, err = tx.Exec("INSERT INTO rubrics_log (question_id,criteria_id,rubrics_id,created_by) VALUES(?,?,?,?)",
				question.Question_id, question.Criteria_ID, question.Rubrics_ID, question.Created_by)
		}

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}