package category

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Category struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	InchargeName string `json:"incharge"`
	MaxTeam      int    `json:"max_team"`
}

type CategoryC struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	InchargeName   string `json:"incharge"`
	MaxTeam        int    `json:"max_team"`
	Registerstatus int    `json:"registration"`
}

type Input struct {
	Id int `json:"id"`
}

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

// to get particular category details
func GetDetail(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var categories CategoryC
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	err = config.Database.QueryRow("SELECT mc.id,mc.category_name,mc.description,mu.name,mc.max_team,mc.registration FROM event_categories ec INNER JOIN m_category mc ON mc.id = ec.category_id INNER JOIN m_events mee ON mee.id = ec.event_id INNER JOIN m_users mu ON mu.id = mc.incharge WHERE ec.status = '1' AND ec.category_id = ?", input.Id).Scan(&categories.Id, &categories.Name, &categories.Description, &categories.InchargeName, &categories.MaxTeam, &categories.Registerstatus)

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

	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}

// to get user name
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Year  string `json:"year"`
}

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

// to insert the team members details
func InsertData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TeamName        string `json:"teamName"`
		EventCategoryID int    `json:"eventCategoryID"`
		User1           string `json:"user1"`
		User2           string `json:"user2"`
		User3           string `json:"user3"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := config.Database.Exec("INSERT INTO event_register (team_name, event_category_id, user_1, user_2, user_3, status, created_by, created_on, updated_on) VALUES (?, ?, ?, ?, ?, '1', ?, NOW(), NOW())",
		req.TeamName, req.EventCategoryID, req.User1, req.User2, req.User3, req.User1)
	fmt.Print(req.User1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}

type Team struct {
	ID               int     `json:"id"`
	TeamName         string  `json:"team_name"`
	EventCategoryID  int     `json:"event_category_id"`
	User1            string  `json:"user_1"`
	User2            *string `json:"user_2"`
	User3            *string `json:"user_3"`
	CreatedBy        string  `json:"created_by"`
	CategoryName     string  `json:"category_name"`
	TeamLeaderName   string  `json:"name1"`
	TeamLeaderMobile string  `json:"phone"`
}

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

// to get team and team details
type TeamI struct {
	ID               int     `json:"id"`
	TeamName         string  `json:"team_name"`
	EventCategoryID  int     `json:"event_category_id"`
	User1            string  `json:"user_1"`
	User1name        string  `json:"name1"`
	User2            *string `json:"user_2"`
	User2name        *string `json:"name2"`
	User3            *string `json:"user_3"`
	User3name        *string `json:"name3"`
	CreatedBy        string  `json:"created_by"`
	CategoryName     string  `json:"category_name"`
	TeamLeaderName   string  `json:"namet"`
	TeamLeaderMobile string  `json:"phone"`
}

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

// to get all events details
type Events struct {
	EventName   string `json:"event_name"`
	Description string `json:"description"`
	Incharge    string `json:"name"`
	EventDate   string `json:"event_date"`
}

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

// To check the team registered or not
type CTeam struct {
	User1     string         `json:"user_1"`
	User1Name string         `json:"user1_name"`
	User2     sql.NullString `json:"user_2"`
	User2Name string         `json:"user2_name"`
	User3     sql.NullString `json:"user_3"`
	User3Name string         `json:"user3_name"`
}

func CheckTeam(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		RollNo string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var isRegistered bool
	rows, err := config.Database.Query(`
        SELECT er.user_1, er.user_2, er.user_3 
        FROM event_register er 
        WHERE er.status='1' 
            AND (er.user_1=? OR er.user_2=? OR er.user_3=?)`,
		requestData.RollNo, requestData.RollNo, requestData.RollNo)

	if err != nil {

		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var team CTeam
		err := rows.Scan(&team.User1, &team.User2, &team.User3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		isRegistered = true
		break
	}
	response := struct {
		IsRegistered bool `json:"isRegistered"`
	}{IsRegistered: isRegistered}
	fmt.Print(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// to get my events data
type MyEvents struct {
	TeamName     string  `json:"team_name"`
	User1        string  `json:"user_1"`
	User1_Name   string  `json:"user_1_name"`
	User2        *string `json:"user_2"`
	User2_Name   *string `json:"user_2_name"`
	User3        *string `json:"user_3"`
	User3_Name   *string `json:"user_3_name"`
	EIncharge    string  `json:"eincharge"`
	CIncharge    string  `json:"cincharge"`
	EventName    string  `json:"event_name"`
	Edesciption  string  `json:"edescription"`
	EventDate    string  `json:"event_date"`
	CategoryName string  `json:"cname"`
	CDescription string  `json:"cdescription"`
	Category_id  int     `json:"event_category_id"`
}

func GetMyEvents(w http.ResponseWriter, r *http.Request) {
	// Parse request 	body
	var requestData struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow(
		"SELECT er.team_name, er.user_1, mu1.name AS user_1_name, er.user_2,mu2.name AS user_2_name, er.user_3,mu3.name AS user_3_name, muu.name AS event_incharge,muc.name AS cincharge, me.event_name, me.description AS edescription,er.event_category_id, mc.category_name AS cname, mc.description AS cdescription, me.event_date FROM `event_register` er INNER JOIN `event_categories` ec ON ec.`id` = er.`event_category_id` INNER JOIN m_category mc ON mc.id=er.event_category_id INNER JOIN m_events me ON me.id = ec.`event_id` INNER JOIN m_users muu ON muu.id=me.incharge LEFT JOIN `m_users` mu1 ON mu1.id = er.`user_1` LEFT JOIN m_users mu2 ON mu2.id = er.`user_2` LEFT JOIN m_users mu3 ON mu3.id = er.`user_3` INNER JOIN m_users muc ON muc.id=mc.incharge WHERE (er.user_1 =? OR er.user_2 =? OR er.user_3 =?)",
		requestData.UserID, requestData.UserID, requestData.UserID)

	var events MyEvents
	err := row.Scan(
		&events.TeamName, &events.User1, &events.User1_Name, &events.User2, &events.User2_Name, &events.User3,
		&events.User3_Name, &events.EIncharge, &events.CIncharge, &events.EventName, &events.Edesciption, &events.Category_id, &events.CategoryName,
		&events.CDescription, &events.EventDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Events not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Events MyEvents `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type GetAllECategory struct {
	CategoryID       int    `json:"id"`
	CategoryName     string `json:"category_name"`
	RegisterCount    int    `json:"category_count"`
	CaDescritpion    string `json:"descritpion"`
	MaxTeams         int    `json:"max_team"`
	CategoryIncharge string `json:"incharge"`
}

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

type users struct {
	RollNo     string `json:"roll_no"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	Year       string `json:"year"`
}

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

type Categorycount struct {
	Ccount int `json:"total_category_count"`
}

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

type RegisterCount struct {
	Rcount int `json:"registercount"`
}

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

type UsersRole struct {
	ID       string `json:"id"`
	UserRole string `json:"user_role"`
}

func GetRole(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT id, role FROM m_users WHERE id = ?", requestData.UserID)
	var events UsersRole
	err := row.Scan(&events.ID, &events.UserRole)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Events not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Events UsersRole `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type UsersRoleC struct {
	ID       string `json:"id"`
	UserRole string `json:"user_role"`
}

func GetRoleC(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT id, role FROM m_users WHERE id = ?", requestData.UserID)
	var events UsersRoleC
	err := row.Scan(&events.ID, &events.UserRole)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Events not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Events UsersRoleC `json:"events"`
	}{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertcategoryData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Category_Name string `json:"category_name"`
		Description   string `json:"description"`
		Max_Team      int    `json:"max_team"`
		Incharge      string `json:"incharge"`
		Created_by    string `json:"created__by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := config.Database.Exec("INSERT INTO m_category (category_name, description, max_team, incharge,status,created_by, created_at, updated_on) VALUES (?, ?, ?, ?, '1', ?, NOW(), NOW())",
		req.Category_Name, req.Description, req.Max_Team, req.Incharge, req.Created_by)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}

type UserCRole struct {
	Name string `json:"name"`
}

func GetCRole(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT name FROM m_users WHERE addcategory_role = '1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []UserCRole
	for rows.Next() {
		var user UserCRole
		if err := rows.Scan(&user.Name); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []UserCRole `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertEventData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EventName   string `json:"event_name"`
		Description string `json:"description"`
		Event_date  string `json:"event_date"`
		Incharge    string `json:"incharge"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	eventDate, err := time.Parse("2006-01-02", req.Event_date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = config.Database.Exec("INSERT INTO m_events (event_name, description, event_date, incharge, status, created_at, updated_on) VALUES (?, ?, ?, ?, '1', NOW(), NOW())",
		req.EventName, req.Description, eventDate, req.Incharge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted successfully",
	}
	function.Response(w, response)
}

type CategoryCountR struct {
	CRcount int `json:"category_count"`
}

type InputR struct {
	Id int `json:"id"`
}

func GetCategoryCountR(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var categories []CategoryCountR
	var temp CategoryCountR

	// Parse the request to get the 'id'
	var input InputR
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}

	row, err := config.Database.Query("SELECT COUNT(*) AS category_count FROM event_register WHERE event_category_id=? AND STATUS='1'", input.Id)

	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}

	for row.Next() {
		err := row.Scan(&temp.CRcount)
		if err != nil {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
			function.Response(w, response)
			return
		}

		tempRow := CategoryCountR{
			CRcount: temp.CRcount,
		}
		categories = append(categories, tempRow)
	}
	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}

type QuestionCategory struct {
	Id   int    `json:"id"`
	Name string `json:"category_name"`
}

func GetCategoryName(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT id,category_name FROM m_category WHERE STATUS = '1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []QuestionCategory
	for rows.Next() {
		var user QuestionCategory
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []QuestionCategory `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Function to get Topics based on category ID
type Topics struct {
	TopicsName string `json:"topics"`
}

type CategoryIdInput struct {
	CId int `json:"id"`
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			response := map[string]interface{}{
				"success": false,
				"error":   "Internal Server Error",
			}
			function.Response(w, response)
		}
	}()

	var response map[string]interface{}
	var categories []Topics
	var input CategoryIdInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	rows, err := config.Database.Query("SELECT topics FROM m_topics WHERE STATUS='1' AND category_id=?", input.CId)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		function.Response(w, response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category Topics
		err := rows.Scan(&category.TopicsName)
		if err != nil {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
			function.Response(w, response)
			return
		}
		categories = append(categories, category)
	}

	response = map[string]interface{}{
		"success": true,
		"data":    categories,
	}
	function.Response(w, response)
}

//function to insert questions

func InsertQuestions(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		CategoryID     int    `json:"category_id"`
		Topics         string `json:"topics"`
		Scenario       string `json:"scenario"`
		Question1      string `json:"question_1"`
		Question_1_Key string `json:"question_1_key"`
		Question2      string `json:"question_2"`
		Question_2_Key string `json:"question_2_key"`
		Question3      string `json:"question_3"`
		Question_3_Key string `json:"question_3_key"`
		Created_by     string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.Database.Begin() // Start a transaction
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Iterate through the questions and insert or update them one by one
	for _, question := range req {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM m_questions WHERE topics = ? AND created_by = ?", question.Topics, question.Created_by).Scan(&count)
		if err != nil {
			tx.Rollback() // Rollback the transaction if there's an error
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			_, err = tx.Exec("UPDATE m_questions SET category_id = ?, scenario = ?, question_1 = ?, question_1_key = ?, question_2 = ?, question_2_key = ?, question_3 = ?, question_3_key = ?, status = '1', updated_on = NOW() WHERE topics = ? AND created_by = ?",
				question.CategoryID, question.Scenario, question.Question1, question.Question_1_Key, question.Question2, question.Question_2_Key, question.Question3, question.Question_3_Key, question.Topics, question.Created_by)
		} else {
			_, err = tx.Exec("INSERT INTO m_questions (category_id,topics,scenario,question_1,question_1_key,question_2,question_2_key,question_3,question_3_key,created_by,status,created_at,updated_on) VALUES (?,?,?,?,?,?,?,?,?,?,'1',NOW(),NOW())",
				question.CategoryID, question.Topics, question.Scenario, question.Question1, question.Question_1_Key, question.Question2, question.Question_2_Key, question.Question3, question.Question_3_Key, question.Created_by)
		}

		if err != nil {
			tx.Rollback() // Rollback the transaction if there's an error
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}


//function to get questions based on user id

type MyQuestions struct {
	Category_Name string `json:"category_name"`
	Topics        string `json:"topics"`
	Scenario      string `json:"scenario"`
	Question_1    string `json:"question_1"`
	Question_2    string `json:"question_2"`
	Question_3    string `json:"question_3"`
}

func GetMyQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		UserID string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := config.Database.Query("SELECT mc.category_name,mq.topics,mq.scenario,mq.question_1,mq.question_2,mq.question_3 FROM m_questions mq INNER JOIN m_category mc ON mc.id=mq.category_id WHERE mq.status='1' AND mq.created_by=?", requestData.UserID)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []MyQuestions
	for rows.Next() {
		var question MyQuestions
		err := rows.Scan(&question.Category_Name, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []MyQuestions `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type TotalQuestion struct {
	Team_Name     string `jons:"team_name"` 
	CreatorName   string `json:"name"`
	Category_Name string `json:"category_name"`
	Topics        string `json:"topics"`
	Scenario      string `json:"scenario"`
	Question_1    string `json:"question_1"`
	Question_2    string `json:"question_2"`
	Question_3    string `json:"question_3"`

}

func TotalQuestions(w http.ResponseWriter, r *http.Request) {

	rows, err := config.Database.Query("SELECT er.team_name,mu.name,mc.category_name,mq.topics,mq.scenario,mq.question_1,mq.question_2,mq.question_3 FROM m_questions mq INNER JOIN m_category mc ON mc.id=mq.category_id INNER JOIN m_users mu ON mu.id=mq.created_by INNER JOIN event_register er ON er.user_1=mq.created_by WHERE mq.status='1'")

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []TotalQuestion
	for rows.Next() {
		var question TotalQuestion
		err := rows.Scan(&question.Team_Name,&question.CreatorName,&question.Category_Name, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2,&question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Questions []TotalQuestion `json:"events"` // Corrected field name here
	}{Questions: questions} // Corrected field name here
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// funciton to get question based on category ID
type GetQuestionResponse struct {
	QuestionID int    `json:"id"`
	Topics     string `json:"topics"`
	Scenario   string `json:"scenario"`
	Question_1 string `json:"question_1"`
	Question_2 string `json:"question_2"`
	Created_by string `json:"created_by"`
}

func GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		Category_ID int    `json:"category_id"`
		Created_by  string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Modify your SQL query to include the OFFSET clause
	rows, err := config.Database.Query("SELECT id,topics, scenario, question_1, question_2, created_by FROM m_questions WHERE category_id=? AND STATUS='1' AND assigned=1 AND created_by!=? LIMIT 10", requestData.Category_ID, requestData.Created_by)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []GetQuestionResponse
	for rows.Next() {
		var question GetQuestionResponse
		err := rows.Scan(&question.QuestionID, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Created_by)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []GetQuestionResponse `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type GetMyCategory struct {
	CategoryID int `json:"event_category_id"`
}

func GetMyCategorys(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		User_ID string `json:"user_1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	row := config.Database.QueryRow("SELECT event_category_id FROM event_register WHERE user_1=? AND STATUS='1'", requestData.User_ID)

	var categoryID int
	err := row.Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No category found for the user", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Prepare response
	response := struct {
		Event GetMyCategory `json:"event"`
	}{Event: GetMyCategory{CategoryID: categoryID}}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// update the assigned id
func UpdateAssignedStatus(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		QuestionIDs []int `json:"id"`
		Assigned    int   `json:"assigned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if len(requestData.QuestionIDs) == 0 {
		http.Error(w, "No question IDs provided", http.StatusBadRequest)
		return
	}

	questionIDString := ""
	for i, id := range requestData.QuestionIDs {
		if i > 0 {
			questionIDString += ","
		}
		questionIDString += fmt.Sprintf("%d", id)
	}

	_, err := config.Database.Exec(
		"UPDATE m_questions SET assigned=? WHERE id IN ("+questionIDString+")",
		requestData.Assigned)

	if err != nil {
		http.Error(w, "Failed to update assigned status in the database", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	response := struct {
		Message string `json:"message"`
	}{Message: "Assigned status updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//function to insert the assigned questions

type QuestionAssigned struct {
	CategoryID int    `json:"category_id"`
	QuestionID int    `json:"question_id"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

type RequestData struct {
	CategoryID int    `json:"category_id"`
	QuestionID []int  `json:"question_id"`
	AssignedTo string `json:"assigned_to"`
}

func InsertQuestionAssigned(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists in the database
	var count int
	err := config.Database.QueryRow("SELECT COUNT(*) FROM question_set WHERE assigned_to = ?",
		requestData.AssignedTo).Scan(&count)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	if count > 0 {
		// User already exists, handle accordingly (e.g., update questions)
		// You can put your update logic here
		return // Terminate the function
	}

	// User does not exist, proceed with insertion
	tx, err := config.Database.Begin()
	if err != nil {
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO question_set (category_id, question_id, assigned_to, status) VALUES (?, ?, ?, '1')")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer stmt.Close()

	for _, questionID := range requestData.QuestionID {
		_, err := stmt.Exec(requestData.CategoryID, questionID, requestData.AssignedTo)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error executing statement", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Questions inserted successfully"))
}

type GetMyassignQuestion struct {
	Questions_ID int     `json:"id"`
	Topics       string  `json:"topics"`
	Scenario     *string `json:"scenario"`
	Question_1   *string `json:"question_1"`
	Question_2   *string `json:"question_2"`
	Question_3   *string `json:"question_3"`
}

func GetMyassign(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		User_1 string `json:"user_1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

rows, err := config.Database.Query("SELECT mq.id, mq.topics, mq.scenario, mq.question_1, mq.question_2, mq.question_3 FROM question_set qs INNER JOIN m_questions mq ON mq.id = qs.question_id INNER JOIN event_register er ON er.id = qs.assigned_team_id WHERE er.user_1=?", requestData.User_1)

	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var questions []GetMyassignQuestion
	for rows.Next() {
		var question GetMyassignQuestion
		err := rows.Scan(&question.Questions_ID, &question.Topics, &question.Scenario, &question.Question_1, &question.Question_2, &question.Question_3)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		questions = append(questions, question)
	}

	// Prepare response
	response := struct {
		Events []GetMyassignQuestion `json:"events"`
	}{Events: questions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %s", err)
		return
	}

	for key, values := range r.Form {
		fmt.Fprintf(w, "%s: %s\n", key, values)
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Print("limit")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("saveUrl")
	if err != nil {
		fmt.Print("img")
		http.Error(w, "Error Retrieving the File", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Define the path where you want to save the uploaded images
	// Make sure this path exists on your server
	savePath := "images/" + handler.Filename

	f, err := os.Create(savePath)
	if err != nil {
		fmt.Print("savw")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write the file
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Print("co")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8080/%s", handler.Filename)
	fmt.Fprintf(w, "{\"imageUrl\": \"%s\"}", imageURL)
}

//function insert answers

func InsertAnswers(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		AnsweredBy     string `json:"answered_by"`
		Questionset_ID int    `json:"questionset_id"`
		Question1Ans   string `json:"question_1_ans"`
		Question2Ans   string `json:"question_2_ans"`
		Question3Ans   string `json:"question_3_ans"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.Database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, question := range req {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM m_answers WHERE answered_by=? AND questionset_id=?",
			question.AnsweredBy, question.Questionset_ID).Scan(&count)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			// Data exists, perform an update
			_, err = tx.Exec(`UPDATE m_answers
                              SET question_1_ans=?, question_2_ans=?, question_3_ans=?, updated_on=NOW()
                              WHERE answered_by=? AND questionset_id=?`,
				question.Question1Ans, question.Question2Ans, question.Question3Ans,
				question.AnsweredBy, question.Questionset_ID)
		} else {
			// Data does not exist, perform an insert
			_, err = tx.Exec(`INSERT INTO m_answers
                              (answered_by, questionset_id, question_1_ans, question_2_ans, question_3_ans, status, created_on, updated_on)
                              VALUES (?, ?, ?, ?, ?, '1', NOW(), NOW())`,
				question.AnsweredBy, question.Questionset_ID, question.Question1Ans, question.Question2Ans, question.Question3Ans)
		}

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}

// to get button status

type ActionStatusset struct {
	Id            int `json:"id"`
	Save_question int `json:"save_question"`
	Save_answer   int `json:"save_answer"`
	Save_rubrics  int `json:"save_rubrics"`
}

func ButtonActionStatus(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT id,save_question,save_answer,save_rubrics FROM button_status WHERE STATUS='1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events []ActionStatusset
	for rows.Next() {
		var user ActionStatusset
		if err := rows.Scan(&user.Id, &user.Save_question, &user.Save_answer, &user.Save_rubrics); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, user)
	}
	response := struct {
		Events []ActionStatusset `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//function to insert the rubrics data

func InsertRubricsData(w http.ResponseWriter, r *http.Request) {
	var req []struct {
		Question_id int    `json:"question_id"`
		Criteria_ID int    `json:"criteria_id"`
		Rubrics_ID  int    `json:"selected"`
		Created_by  string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.Database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, question := range req {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM rubrics_log WHERE question_id=? AND criteria_id=? AND created_by=?",
			question.Question_id, question.Criteria_ID, question.Created_by).Scan(&count)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			// Data exists, perform an update
			_, err = tx.Exec("UPDATE rubrics_log SET rubrics_id=? WHERE question_id=? AND criteria_id=? AND created_by=?",
				question.Rubrics_ID, question.Question_id, question.Criteria_ID, question.Created_by)
		} else {
			// Data does not exist, perform an insert
			_, err = tx.Exec("INSERT INTO rubrics_log (question_id,criteria_id,rubrics_id,created_by) VALUES(?,?,?,?)",
				question.Question_id, question.Criteria_ID, question.Rubrics_ID, question.Created_by)
		}

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Data inserted/updated successfully",
	}
	function.Response(w, response)
}
