package roles

type UsersRole struct {
	ID       string `json:"id"`
	UserRole string `json:"user_role"`
}

type UsersRoleC struct {
	ID       string `json:"id"`
	UserRole string `json:"user_role"`
}

type UserCRole struct {
	Name string `json:"name"`
}

type CategoryCountR struct {
	CRcount int `json:"category_count"`
}

type InputR struct {
	Id int `json:"id"`
}