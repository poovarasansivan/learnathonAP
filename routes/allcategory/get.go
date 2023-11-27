package allcategory

import (
	"encoding/json"
	"learnathon/config"
	"net/http"
)

func GetAllEVCategory(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT mc.id, mc.category_name, mc.description AS description, mc.max_team, mc.incharge, COUNT(er.event_category_id) AS category_count FROM m_category mc LEFT JOIN event_register er ON er.event_category_id = mc.id WHERE mc.status = '1' GROUP BY mc.id;	")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []GetAllECategory
	for rows.Next() {
		var team GetAllECategory

		err := rows.Scan(&team.CategoryID, &team.CategoryName, &team.CaDescritpion, &team.MaxTeams, &team.CategoryIncharge, &team.RegisterCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams = append(teams, team)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}