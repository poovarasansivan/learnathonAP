package events
import (
	"database/sql"
	"learnathon/config"
	"learnathon/function"
	"net/http"
)

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var eventdata []Events
	var temp Events

	row, err := config.Database.Query("SELECT me.event_name,me.description,me.event_date,mu.name FROM m_events me INNER JOIN m_users mu ON mu.id = me.incharge WHERE me.status ='1'")

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
		err := row.Scan(&temp.EventName, &temp.Description, &temp.EventDate, &temp.Incharge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tempRow := Events{
			EventName:   temp.EventName,
			Description: temp.Description,
			Incharge:    temp.Incharge,
			EventDate:   temp.EventDate,
		}
		eventdata = append(eventdata, tempRow)
	}

	response = map[string]interface{}{
		"success": true,
		"data":    eventdata,
	}
	function.Response(w, response)
}