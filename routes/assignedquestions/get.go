package assignedquestions

import (
	"encoding/json"
	"fmt"
	"learnathon/config"
	"log"
	"net/http"
)

func UpdateAssignedStatus(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		QuestionIDs []int `json:"id"`
		Assigned    int   `json:"assigned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if len(requestData.QuestionIDs) == 0 {
		http.Error(w, "No question IDs provided", http.StatusBadRequest)
		return
	}

	questionIDString := ""
	for i, id := range requestData.QuestionIDs {
		if i > 0 {
			questionIDString += ","
		}
		questionIDString += fmt.Sprintf("%d", id)
	}

	_, err := config.Database.Exec(
		"UPDATE m_questions SET assigned=? WHERE id IN ("+questionIDString+")",
		requestData.Assigned)

	if err != nil {
		http.Error(w, "Failed to update assigned status in the database", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	response := struct {
		Message string `json:"message"`
	}{Message: "Assigned status updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertQuestionAssigned(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists in the database
	var count int
	err := config.Database.QueryRow("SELECT COUNT(*) FROM question_set WHERE assigned_to = ?",
		requestData.AssignedTo).Scan(&count)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	if count > 0 {
		// User already exists, handle accordingly (e.g., update questions)
		// You can put your update logic here
		return // Terminate the function
	}

	// User does not exist, proceed with insertion
	tx, err := config.Database.Begin()
	if err != nil {
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO question_set (category_id, question_id, assigned_to, status) VALUES (?, ?, ?, '1')")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer stmt.Close()

	for _, questionID := range requestData.QuestionID {
		_, err := stmt.Exec(requestData.CategoryID, questionID, requestData.AssignedTo)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error executing statement", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Questions inserted successfully"))
}



func GetMyassign(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		User_1 string `json:"user_1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := config.Database.Query("SELECT mq.id,mq.`topics`,mq.`scenario`,mq.`question_1`,mq.`question_2`,mq.`question_3` FROM question_set qs INNER JOIN m_questions mq ON mq.`id`=qs.`question_id` INNER JOIN event_register er ON er.`id`=qs.`assigned_team_id` WHERE er.`user_1`=?", requestData.User_1)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []GetMyassignQuestion
	for rows.Next() {
		var question GetMyassignQuestion
		err := rows.Scan(&question.Questions_ID, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []GetMyassignQuestion `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}