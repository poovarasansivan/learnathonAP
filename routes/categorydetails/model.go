package categorydetails

type CategoryDetails struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	InchargeName   string `json:"incharge"`
	MaxTeam        int    `json:"max_team"`
	Registerstatus int    `json:"registration"`
}

type Input struct {
	Id int `json:"id"`
}