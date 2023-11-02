package category

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
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
	var categories Category
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

	err = config.Database.QueryRow("SELECT mc.id,mc.category_name,mc.description,mu.name,mc.max_team,mc.registration FROM event_categories ec INNER JOIN m_category mc ON mc.id = ec.category_id INNER JOIN m_events mee ON mee.id = ec.event_id INNER JOIN m_users mu ON mu.id = mc.incharge WHERE ec.status = '1' AND ec.category_id = ?", input.Id).Scan(&categories.Id, &categories.Name, &categories.Description, &categories.InchargeName, &categories.MaxTeam)

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
	ID               int            `json:"id"`
	TeamName         string         `json:"team_name"`
	EventCategoryID  int            `json:"event_category_id"`
	User1            string         `json:"user_1"`
	User2            sql.NullString `json:"user_2"`
	User3            sql.NullString `json:"user_3"`
	CreatedBy        string         `json:"created_by"`
	CategoryName     string         `json:"category_name"`
	TeamLeaderName   string         `json:"name1"`
	TeamLeaderMobile int            `json:"phone"`
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	rows, err := config.Database.Query("SELECT er.user_3,er.user_2,er.user_1,er.team_name,er.id,er.event_category_id,er.created_by,mc.category_name,mu.name AS name1,mu.phone FROM event_register er INNER JOIN m_category mc ON mc.id=er.event_category_id INNER JOIN m_users mu ON mu.id=er.user_1 WHERE er.status='1'")

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
	ID               int            `json:"id"`
	TeamName         string         `json:"team_name"`
	EventCategoryID  int            `json:"event_category_id"`
	User1            string         `json:"user_1"`
	User1name        string         `json:"name1"`
	User2            sql.NullString `json:"user_2"`
	User2name        string         `json:"name2"`
	User3            sql.NullString `json:"user_3"`
	User3name        string         `json:"name3"`
	CreatedBy        string         `json:"created_by"`
	CategoryName     string         `json:"category_name"`
	TeamLeaderName   string         `json:"namet"`
	TeamLeaderMobile int            `json:"phone"`
}

func GetTeamByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	teamID := params["team_id"]

	query := `
    SELECT 
    er.user_3,
    mu3.name AS name3,
    er.user_2,
    mu2.name AS name2,
    er.user_1,
    mu.name AS name1,
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var team1 TeamI
	for rows.Next() {
		err := rows.Scan(&team1.User1, &team1.User1name, &team1.User2, &team1.User2name, &team1.User3, &team1.User3name, &team1.TeamName, &team1.ID, &team1.EventCategoryID, &team1.CreatedBy, &team1.CategoryName, &team1.TeamLeaderName, &team1.TeamLeaderMobile)
		if err != nil {
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
	TeamName     string `json:"team_name"`
	User1        string `json:"user_1"`
	User1_Name   string `json:"user_1_name"`
	User2        string `json:"user_2"`
	User2_Name   string `json:"user_2_name"`
	User3        string `json:"user_3"`
	User3_Name   string `json:"user_3_name"`
	EIncharge    string `json:"eincharge"`
	CIncharge    string `json:"cincharge"`
	EventName    string `json:"event_name"`
	Edesciption  string `json:"edescription"`
	EventDate    string `json:"event_date"`
	CategoryName string `json:"cname"`
	CDescription string `json:"cdescription"`
	Category_id  int    `json:"event_category_id"`
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
		"SELECT er.team_name, er.user_1, mu1.name AS user_1_name, er.user_2,mu2.name AS user_2_name, er.user_3,mu3.name AS user_3_name, muu.name AS event_incharge,muc.name AS cincharge, me.event_name, me.description AS edescription,er.event_category_id, mc.category_name AS cname, mc.description AS cdescription, me.event_date FROM `event_register` er INNER JOIN `event_categories` ec ON ec.`id` = er.`event_category_id` INNER JOIN m_category mc ON mc.id=er.event_category_id INNER JOIN m_events me ON me.id = ec.`event_id` INNER JOIN m_users muu ON muu.id=me.incharge INNER JOIN `m_users` mu1 ON mu1.id = er.`user_1` INNER JOIN m_users mu2 ON mu2.id = er.`user_2` INNER JOIN m_users mu3 ON mu3.id = er.`user_3` INNER JOIN m_users muc ON muc.id=mc.incharge WHERE (er.user_1 = ? OR er.user_2 =? OR er.user_3 = ?)",
		requestData.UserID, requestData.UserID, requestData.UserID)
	
	var events MyEvents
	err := row.Scan(
		&events.TeamName, &events.User1, &events.User1_Name, &events.User2, &events.User2_Name, &events.User3,
		&events.User3_Name, &events.EIncharge,&events.CIncharge,&events.EventName, &events.Edesciption,&events.Category_id, &events.CategoryName,
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
