package teamdetails
type TeamI struct {
	ID               int     `json:"id"`
	TeamName         string  `json:"team_name"`
	EventCategoryID  int     `json:"event_category_id"`
	User1            string  `json:"user_1"`
	User1name        string  `json:"name1"`
	User2            *string `json:"user_2"`
	User2name        *string `json:"name2"`
	User3            *string `json:"user_3"`
	User3name        *string `json:"name3"`
	CreatedBy        string  `json:"created_by"`
	CategoryName     string  `json:"category_name"`
	TeamLeaderName   string  `json:"namet"`
	TeamLeaderMobile string  `json:"phone"`
}