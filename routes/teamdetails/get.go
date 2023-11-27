package teamdetails

import (
	"encoding/json"
	"fmt"
	"learnathon/config"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTeamByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	teamID := params["team_id"]

	query := `
    SELECT 
    er.user_1,
    mu1.name AS name1,
    er.user_2,
    mu2.name AS name2,
    er.user_3,
    mu3.name AS name3,
    er.team_name,
    er.id,
    er.event_category_id,
    er.created_by,
	mc.category_name,
    mu1.name AS namet,
    mu.phone
FROM 
    event_register er 
INNER JOIN 
    m_category mc ON mc.id = er.event_category_id 
INNER JOIN 
    m_users mu ON mu.id = er.user_1 
LEFT JOIN 
    m_users mu1 ON mu1.id = er.user_1
LEFT JOIN 
    m_users mu2 ON mu2.id = er.user_2
LEFT JOIN 
    m_users mu3 ON mu3.id = er.user_3
WHERE 
    er.id = ? AND er.status = '1'

`

	rows, err := config.Database.Query(query, teamID)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var team1 TeamI
	for rows.Next() {
		err := rows.Scan(&team1.User1, &team1.User1name, &team1.User2, &team1.User2name, &team1.User3, &team1.User3name, &team1.TeamName, &team1.ID, &team1.EventCategoryID, &team1.CreatedBy, &team1.CategoryName, &team1.TeamLeaderName, &team1.TeamLeaderMobile)
		if err != nil {
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(team1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}