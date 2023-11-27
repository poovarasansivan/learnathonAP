package rubrics

type Rubrics struct {
	RubricsID   string `json:"rubrics_id"`
	RubricsName string `json:"rubrics_name"`
}

type Criteria struct {
	CriteriaID   int       `json:"criteria_id"`
	CriteriaName string    `json:"criteria_name"`
	Rubrics      []Rubrics `json:"rubrics"`
}
