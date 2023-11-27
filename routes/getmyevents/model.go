package getmyevents

type MyEvents struct {
	TeamName     string  `json:"team_name"`
	User1        string  `json:"user_1"`
	User1_Name   string  `json:"user_1_name"`
	User2        *string `json:"user_2"`
	User2_Name   *string `json:"user_2_name"`
	User3        *string `json:"user_3"`
	User3_Name   *string `json:"user_3_name"`
	EIncharge    string  `json:"eincharge"`
	CIncharge    string  `json:"cincharge"`
	EventName    string  `json:"event_name"`
	Edesciption  string  `json:"edescription"`
	EventDate    string  `json:"event_date"`
	CategoryName string  `json:"cname"`
	CDescription string  `json:"cdescription"`
	Category_id  int     `json:"event_category_id"`
}

type GetMyCategory struct {
	CategoryID int `json:"event_category_id"`
}