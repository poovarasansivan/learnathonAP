package answers
import (
	"encoding/json"
	"learnathon/config"
	"learnathon/function"

	"net/http"
)


func InsertAnswers(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		AnsweredBy     string `json:"answered_by"`
		Questionset_ID int    `json:"questionset_id"`
		Question1Ans   string `json:"question_1_ans"`
		Question2Ans   string `json:"question_2_ans"`
		Question3Ans   string `json:"question_3_ans"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
		err := tx.QueryRow("SELECT COUNT(*) FROM m_answers WHERE answered_by=? AND questionset_id=?",
			question.AnsweredBy, question.Questionset_ID).Scan(&count)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			// Data exists, perform an update
			_, err = tx.Exec(`UPDATE m_answers
                              SET question_1_ans=?, question_2_ans=?, question_3_ans=?, updated_on=NOW()
                              WHERE answered_by=? AND questionset_id=?`,
				question.Question1Ans, question.Question2Ans, question.Question3Ans,
				question.AnsweredBy, question.Questionset_ID)
		} else {
			// Data does not exist, perform an insert
			_, err = tx.Exec(`INSERT INTO m_answers
                              (answered_by, questionset_id, question_1_ans, question_2_ans, question_3_ans, status, created_on, updated_on)
                              VALUES (?, ?, ?, ?, ?, '1', NOW(), NOW())`,
				question.AnsweredBy, question.Questionset_ID, question.Question1Ans, question.Question2Ans, question.Question3Ans)
		}

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}
