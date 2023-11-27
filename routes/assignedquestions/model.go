package assignedquestions

type QuestionAssigned struct {
	CategoryID int    `json:"category_id"`
	QuestionID int    `json:"question_id"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

type RequestData struct {
	CategoryID int    `json:"category_id"`
	QuestionID []int  `json:"question_id"`
	AssignedTo string `json:"assigned_to"`
}
type GetMyassignQuestion struct {
	Questions_ID int     `json:"id"`
	Topics       string  `json:"topics"`
	Scenario     *string `json:"scenario"`
	Question_1   *string `json:"question_1"`
	Question_2   *string `json:"question_2"`
	Question_3   *string `json:"question_3"`
}