package questions

import (
	"encoding/json"
	"fmt"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func InsertQuestions(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		CategoryID     int    `json:"category_id"`
		Topics         string `json:"topics"`
		Scenario       string `json:"scenario"`
		Question1      string `json:"question_1"`
		Question_1_Key string `json:"question_1_key"`
		Question2      string `json:"question_2"`
		Question_2_Key string `json:"question_2_key"`
		Question3      string `json:"question_3"`
		Question_3_Key string `json:"question_3_key"`
		Created_by     string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.Database.Begin() // Start a transaction
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Iterate through the questions and insert or update them one by one
	for _, question := range req {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM m_questions WHERE topics = ? AND created_by = ?", question.Topics, question.Created_by).Scan(&count)
		if err != nil {
			tx.Rollback() // Rollback the transaction if there's an error
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			_, err = tx.Exec("UPDATE m_questions SET category_id = ?, scenario = ?, question_1 = ?, question_1_key = ?, question_2 = ?, question_2_key = ?, question_3 = ?, question_3_key = ?, status = '1', updated_on = NOW() WHERE topics = ? AND created_by = ?",
				question.CategoryID, question.Scenario, question.Question1, question.Question_1_Key, question.Question2, question.Question_2_Key, question.Question3, question.Question_3_Key, question.Topics, question.Created_by)
		} else {
			_, err = tx.Exec("INSERT INTO m_questions (category_id,topics,scenario,question_1,question_1_key,question_2,question_2_key,question_3,question_3_key,created_by,status,created_at,updated_on) VALUES (?,?,?,?,?,?,?,?,?,?,'1',NOW(),NOW())",
				question.CategoryID, question.Topics, question.Scenario, question.Question1, question.Question_1_Key, question.Question2, question.Question_2_Key, question.Question3, question.Question_3_Key, question.Created_by)
		}

		if err != nil {
			tx.Rollback() // Rollback the transaction if there's an error
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}

func GetMyQuestions(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := config.Database.Query("SELECT mc.category_name,mq.topics,mq.scenario,mq.question_1,mq.question_2,mq.question_3 FROM m_questions mq INNER JOIN m_category mc ON mc.id=mq.category_id WHERE mq.status='1' AND mq.created_by=?", requestData.UserID)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []MyQuestions
	for rows.Next() {
		var question MyQuestions
		err := rows.Scan(&question.Category_Name, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	response := struct {
		Events []MyQuestions `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TotalQuestions(w http.ResponseWriter, r *http.Request) {

	rows, err := config.Database.Query("SELECT er.team_name,mu.name,mc.category_name,mq.topics,mq.scenario,mq.question_1,mq.question_2,mq.question_3 FROM m_questions mq INNER JOIN m_category mc ON mc.id=mq.category_id INNER JOIN m_users mu ON mu.id=mq.created_by INNER JOIN event_register er ON er.user_1=mq.created_by WHERE mq.status='1'")

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []TotalQuestion
	for rows.Next() {
		var question TotalQuestion
		err := rows.Scan(&question.Team_Name, &question.CreatorName, &question.Category_Name, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Questions []TotalQuestion `json:"events"` // Corrected field name here
	}{Questions: questions} // Corrected field name here
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		Category_ID int    `json:"category_id"`
		Created_by  string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Modify your SQL query to include the OFFSET clause
	rows, err := config.Database.Query("SELECT id,topics, scenario, question_1, question_2, created_by FROM m_questions WHERE category_id=? AND STATUS='1' AND assigned=1 AND created_by!=? LIMIT 10", requestData.Category_ID, requestData.Created_by)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []GetQuestionResponse
	for rows.Next() {
		var question GetQuestionResponse
		err := rows.Scan(&question.QuestionID, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Created_by)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []GetQuestionResponse `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}