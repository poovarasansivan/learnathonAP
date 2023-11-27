package overallusers
import (
	"encoding/json"
	"learnathon/config"
	"net/http"
)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT mu.id AS roll_no,mu.name,mu.email,mu.phone,mu.department,mu.year FROM m_users mu WHERE mu.status='1'")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []users
	for rows.Next() {
		var team users

		err := rows.Scan(&team.RollNo, &team.Name, &team.Email, &team.Phone, &team.Department, &team.Year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams = append(teams, team)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
