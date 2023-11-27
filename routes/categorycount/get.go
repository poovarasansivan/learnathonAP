package categorycount

import (
	"encoding/json"
	"learnathon/config"
	"net/http"
)

func GetCcount(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT COUNT(id) AS total_category_count FROM m_category WHERE STATUS = '1';")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []Categorycount
	for rows.Next() {
		var team Categorycount

		err := rows.Scan(&team.Ccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams = append(teams, team)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}