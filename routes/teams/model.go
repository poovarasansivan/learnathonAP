package teams

type Team struct {
	ID               int     `json:"id"`
	TeamName         string  `json:"team_name"`
	EventCategoryID  int     `json:"event_category_id"`
	User1            string  `json:"user_1"`
	User2            *string `json:"user_2"`
	User3            *string `json:"user_3"`
	CreatedBy        string  `json:"created_by"`
	CategoryName     string  `json:"category_name"`
	TeamLeaderName   string  `json:"name1"`
	TeamLeaderMobile string  `json:"phone"`
}
