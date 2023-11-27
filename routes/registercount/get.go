package registercount

import (
	"encoding/json"
	"learnathon/config"
	"net/http"
)

func GetRegisterCount(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT COUNT(id) AS registercount FROM event_register WHERE STATUS = '1';")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []RegisterCount
	for rows.Next() {
		var team RegisterCount

		err := rows.Scan(&team.Rcount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams = append(teams, team)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}