package mcq

import (
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func McqQuestions(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		Question   string `json:"question"`
		Option1    string `json:"option1"`
		Option2    string `json:"option2"`
		Option3    string `json:"option3"`
		Option4    string `json:"option4"`
		Answer     string `json:"correct_ans"`
		Created_by string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, question := range req {
		_, err := config.Database.Exec("INSERT INTO mcq_questions (question, option1, option2, option3, option4, correct_ans, created_by,created_on) VALUES (?, ?, ?, ?, ?, ?,?, NOW())",
			question.Question, question.Option1, question.Option2, question.Option3, question.Option4, question.Answer, question.Created_by)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}

func McqEvalution(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT id,question,option1,option2,option3,option4,correct_ans FROM mcq_questions WHERE STATUS='1'")

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var mcqquestions []Mcqevalution
	for rows.Next() {
		var question Mcqevalution
		err := rows.Scan(&question.Id, &question.Question, &question.Option_1, &question.Option_2, &question.Option_3, &question.Option_4, &question.Correct_ans)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		mcqquestions = append(mcqquestions, question)
	}

	// Prepare response
	response := struct {
		Questions []Mcqevalution `json:"events"` // Corrected field name here
	}{Questions: mcqquestions} // Corrected field name here
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Mymcqquestions(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := config.Database.Query("SELECT id,question,option1,option2,option3,option4,correct_ans FROM mcq_questions WHERE STATUS='1' AND created_by=?", requestData.UserID)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []MyMcq
	for rows.Next() {
		var question MyMcq
		err := rows.Scan(&question.Id, &question.Question, &question.Option_1, &question.Option_2, &question.Option_3, &question.Option_4,&question.Correct_ans)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	response := struct {
		Events []MyMcq `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


