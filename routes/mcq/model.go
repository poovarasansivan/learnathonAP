package mcq

type Mcqevalution struct {
	Id          int    `json:"id"`
	Question    string `json:"question"`
	Option_1    string `json:"option1"`
	Option_2    string `json:"option2"`
	Option_3    string `json:"option3"`
	Option_4    string `json:"option4"`
	Correct_ans string `json:"correct_ans"`
}

type MyMcq struct {
	Id          int    `json:"id"`
	Question    string `json:"question"`
	Option_1    string `json:"option1"`
	Option_2    string `json:"option2"`
	Option_3    string `json:"option3"`
	Option_4    string `json:"option4"`
	Correct_ans string `json:"correct_ans"`
}
