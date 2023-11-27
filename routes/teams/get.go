package teams

import (
	"encoding/json"

	"learnathon/config"
	"net/http"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT er.user_1,er.user_2,er.user_3,er.team_name,er.id,er.event_category_id,er.created_by,mc.category_name,mu.name AS name1,mu.phone FROM event_register er INNER JOIN m_category mc ON mc.id=er.event_category_id INNER JOIN m_users mu ON mu.id=er.user_1 WHERE er.status='1'")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team

		err := rows.Scan(&team.User1, &team.User2, &team.User3, &team.TeamName, &team.ID, &team.EventCategoryID, &team.CreatedBy, &team.CategoryName, &team.TeamLeaderName, &team.TeamLeaderMobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams = append(teams, team)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
