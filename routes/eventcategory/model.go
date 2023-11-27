package eventcategory

type Category struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	InchargeName string `json:"incharge"`
	MaxTeam      int    `json:"max_team"`
}