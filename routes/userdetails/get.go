package userdetails

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserByName(w http.ResponseWriter, r *http.Request) {
	rollno := mux.Vars(r)["rollno"]

	row := config.Database.QueryRow("SELECT id, name, year, email FROM m_users WHERE id=?", rollno)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Year, &user.Email)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}