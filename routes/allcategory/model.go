package allcategory

type GetAllECategory struct {
	CategoryID       int    `json:"id"`
	CategoryName     string `json:"category_name"`
	RegisterCount    int    `json:"category_count"`
	CaDescritpion    string `json:"descritpion"`
	MaxTeams         int    `json:"max_team"`
	CategoryIncharge string `json:"incharge"`
}