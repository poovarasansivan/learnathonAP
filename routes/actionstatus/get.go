package actionstatus

import (
	"encoding/json"
	"fmt"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func ButtonActionStatus(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT id,save_question,save_answer,registration,save_mcq,save_rubrics FROM button_status WHERE STATUS='1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []ActionStatusset
	for rows.Next() {
		var user ActionStatusset
		if err := rows.Scan(&user.Id, &user.Save_question, &user.Save_answer,&user.Registration,&user.Save_mcq, &user.Save_rubrics); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []ActionStatusset `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Questionstatus(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		User_id     string  `json:"user_id"`
		Question_id int    `json:"question_id"`
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
		err := tx.QueryRow("SELECT COUNT(*) FROM question_status WHERE question_id=? and user_id=?",
		question.Question_id,question.User_id).Scan(&count)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			// Data exists, perform an update
			_, err = tx.Exec("UPDATE question_status SET STATUS='2' WHERE question_id=? and user_id=?",question.Question_id,question.User_id)
		} else {
			// Data does not exist, perform an insert
			_, err = tx.Exec("INSERT INTO question_status (user_id,question_id,STATUS) VALUES(?,?,'2')",question.User_id,question.Question_id)
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

func GetQuestionSubmitstatus(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		User_id string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := config.Database.Query("SELECT question_id,STATUS FROM question_status WHERE user_id=?", requestData.User_id)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []GetQuestionStatus
	for rows.Next() {
		var question GetQuestionStatus
		err := rows.Scan(&question.Questions_ID,&question.Status)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []GetQuestionStatus `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}