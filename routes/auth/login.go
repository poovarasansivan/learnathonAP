package auth

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)

type UserInput struct {
	Email string `json:"email"`
}

type UserAccount struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	var response map[string]interface{}
	var data UserAccount

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	err = config.Database.QueryRow("SELECT id,NAME,email,phone FROM m_users WHERE email= ? AND STATUS ='1'", input.Email).Scan(&data.UserId, &data.Name, &data.Email, &data.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No User Found",
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

	response = map[string]interface{}{
		"success": true,
		"user":    data,
	}
	function.Response(w, response)
}
