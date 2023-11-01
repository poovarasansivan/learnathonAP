package main

import (
	"fmt"
	"learnathon/config"
	"learnathon/routes"
	"learnathon/routes/auth"
	"learnathon/routes/category"

	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	config.ConnectDB()
	defer config.Database.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/", routes.Sample).Methods("POST")
	router.HandleFunc("/auth/login", auth.Login).Methods("POST")

	router.HandleFunc("/api/category/getAll", category.GetAllCategory).Methods("GET")
	router.HandleFunc("/api/category/getDetails", category.GetDetail).Methods("POST")
	router.HandleFunc("/api/users/{rollno}", category.GetUserByName).Methods("GET")
	router.HandleFunc("/api/insertData", category.InsertData).Methods("POST")
	router.HandleFunc("/api/users", category.GetUsers).Methods("GET")
	router.HandleFunc("/api/teams", category.GetTeams).Methods("GET")
	router.HandleFunc("/api/category/GetEvents", category.GetAllEvents).Methods("GET")
	router.HandleFunc("/api/teamsid/{team_id}", category.GetTeamByID).Methods("GET")
	router.HandleFunc("/api/CheckTeam", category.CheckTeam).Methods("POST")
	router.HandleFunc("/api/GetMyEvents", category.GetMyEvents).Methods("POST")
	router.HandleFunc("/api/GetEVCategory", category.GetAllEVCategory).Methods("GET")
	router.HandleFunc("/api/GetCcount", category.GetCcount).Methods("GET")
	router.HandleFunc("/api/GetRcount", category.GetRegisterCount).Methods("GET")
	router.HandleFunc("/api/GetUserRole", category.GetRole).Methods("POST")
	router.HandleFunc("/api/GetUserRoleC", category.GetRoleC).Methods("POST")
	router.HandleFunc("/api/Insertcategory", category.InsertcategoryData).Methods("POST")
	router.HandleFunc("/api/GetUserAdd", category.GetCRole).Methods("POST")
	router.HandleFunc("/api/AddEvents", category.InsertEventData).Methods("POST")
	router.HandleFunc("/api/GetCategoryC", category.GetCategoryCountR).Methods("POST")

	c := cors.AllowAll()

	fmt.Print("Running....")
	handler := c.Handler(router)
	http.Handle("/", handlers.LoggingHandler(os.Stdout, handler))

	http.ListenAndServe(":8080", nil)

}
