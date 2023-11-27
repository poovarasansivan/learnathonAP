package questions

type MyQuestions struct {
	Category_Name string `json:"category_name"`
	Topics        string `json:"topics"`
	Scenario      string `json:"scenario"`
	Question_1    string `json:"question_1"`
	Question_2    string `json:"question_2"`
	Question_3    string `json:"question_3"`
}

type TotalQuestion struct {
	Team_Name     string `jons:"team_name"`
	CreatorName   string `json:"name"`
	Category_Name string `json:"category_name"`
	Topics        string `json:"topics"`
	Scenario      string `json:"scenario"`
	Question_1    string `json:"question_1"`
	Question_2    string `json:"question_2"`
	Question_3    string `json:"question_3"`
}

type GetQuestionResponse struct {
	QuestionID int    `json:"id"`
	Topics     string `json:"topics"`
	Scenario   string `json:"scenario"`
	Question_1 string `json:"question_1"`
	Question_2 string `json:"question_2"`
	Created_by string `json:"created_by"`
}