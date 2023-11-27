package eventcategory

import (
	"database/sql"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)
func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var categories []Category
	var temp Category

	row, err := config.Database.Query("SELECT mc.id,mc.category_name,mc.description,mu.name,mc.max_team FROM event_categories ec INNER JOIN m_category mc ON mc.id = ec.category_id INNER JOIN m_users mu ON mu.id = mc.incharge WHERE ec.status = '1'")

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No Request",
			}
		} else {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		function.Response(w, response)
		return
	}

	for row.Next() {
		err := row.Scan(&temp.Id, &temp.Name, &temp.Description, &temp.InchargeName, &temp.MaxTeam)
		if err != nil {
			panic(err.Error)
		}

		tempRow := Category{
			Id:           temp.Id,
			Name:         temp.Name,
			Description:  temp.Description,
			InchargeName: temp.InchargeName,
			MaxTeam:      temp.MaxTeam,
		}
		categories = append(categories, tempRow)
	}
	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}
