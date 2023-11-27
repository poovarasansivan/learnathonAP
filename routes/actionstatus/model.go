package actionstatus

type ActionStatusset struct {
	Id            int `json:"id"`
	Save_question int `json:"save_question"`
	Save_answer   int `json:"save_answer"`
	Registration  int `json:"registration"`
	Save_mcq      int `json:"save_mcq"`
	Save_rubrics  int `json:"save_rubrics"`
}

type GetQuestionStatus struct {
	Questions_ID int `json:"question_id"`
	Status       int `json:"status"`
}
