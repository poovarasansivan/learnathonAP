package main

import (
	"fmt"
	"learnathon/config"
	"learnathon/routes"
	"learnathon/routes/actionstatus"
	"learnathon/routes/allcategory"
	"learnathon/routes/answers"
	"learnathon/routes/assignedquestions"
	"learnathon/routes/auth"
	"learnathon/routes/mcq"

	"learnathon/routes/categorycount"
	"learnathon/routes/categorydata"
	"learnathon/routes/categorydetails"
	"learnathon/routes/eventcategory"
	"learnathon/routes/events"
	"learnathon/routes/getmyevents"
	"learnathon/routes/image"
	"learnathon/routes/inserteventdata"
	"learnathon/routes/overallusers"
	"learnathon/routes/questions"
	"learnathon/routes/registercount"
	"learnathon/routes/registerdata"
	"learnathon/routes/roles"
	"learnathon/routes/rubrics"
	"learnathon/routes/teamdetails"
	"learnathon/routes/teams"
	"learnathon/routes/topics"
	"learnathon/routes/userdetails"

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
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")

	router.HandleFunc("/api/getAll", eventcategory.GetAllCategory).Methods("GET")
	router.HandleFunc("/api/getDetails", categorydetails.GetDetail).Methods("POST")
	router.HandleFunc("/api/userdetails/{rollno}", userdetails.GetUserByName).Methods("GET")
    router.HandleFunc("/api/insertData", registerdata.InsertData).Methods("POST")
	router.HandleFunc("/api/users", overallusers.GetUsers).Methods("GET")
	router.HandleFunc("/api/teams", teams.GetTeams).Methods("GET")
	router.HandleFunc("/api/GetEvents", events.GetAllEvents).Methods("GET")
	router.HandleFunc("/api/teamsid/{team_id}", teamdetails.GetTeamByID).Methods("GET")
	router.HandleFunc("/api/CheckTeam", registerdata.CheckTeam).Methods("POST")
	router.HandleFunc("/api/GetMyEvents", getmyevents.GetMyEvents).Methods("POST")
	router.HandleFunc("/api/GetEVCategory", allcategory.GetAllEVCategory).Methods("GET")
    router.HandleFunc("/api/GetCcount", categorycount.GetCcount).Methods("GET")
	router.HandleFunc("/api/GetRcount", registercount.GetRegisterCount).Methods("GET")
	router.HandleFunc("/api/GetUserRole", roles.GetRole).Methods("POST")
	router.HandleFunc("/api/GetUserRoleC", roles.GetRoleC).Methods("POST")	
	router.HandleFunc("/api/Insertcategory", categorydata.InsertcategoryData).Methods("POST")
	router.HandleFunc("/api/GetUserAdd", roles.GetCRole).Methods("POST")
	router.HandleFunc("/api/AddEvents", inserteventdata.InsertEventData).Methods("POST")
	router.HandleFunc("/api/GetCategoryC", roles.GetCategoryCountR)
    router.HandleFunc("/api/GetCName", categorydata.GetCategoryName).Methods("GET")
    router.HandleFunc("/api/GetTopics", topics.GetTopics).Methods("POST")
	router.HandleFunc("/api/insertQuestion", questions.InsertQuestions).Methods("POST")
	router.HandleFunc("/api/GetMyQuestion", questions.GetMyQuestions).Methods("POST")
	router.HandleFunc("/api/TotalQuestion", questions.TotalQuestions).Methods("GET")
	router.HandleFunc("/api/GetAllQuestion", questions.GetAllQuestions).Methods("POST")
	router.HandleFunc("/api/GetMyCategory", getmyevents.GetMyCategorys).Methods("POST")
	router.HandleFunc("/api/updateAssigned", assignedquestions.UpdateAssignedStatus).Methods("POST")
	router.HandleFunc("/api/InsertAssignQuestion", assignedquestions.InsertQuestionAssigned).Methods("POST")
	router.HandleFunc("/api/GetMyassignQuestions", assignedquestions.GetMyassign).Methods("POST")
    router.HandleFunc("/api/InsertAnswer", answers.InsertAnswers).Methods("POST")
	router.HandleFunc("/api/ButtonStatus", actionstatus.ButtonActionStatus).Methods("GET")
	router.HandleFunc("/api/RubricsData", rubrics.InsertRubricsData).Methods("POST")
	router.HandleFunc("/api/QuestionStatus", actionstatus.Questionstatus).Methods("POST")
	router.HandleFunc("/api/QuestionSubmit", actionstatus.GetQuestionSubmitstatus).Methods("POST")
	router.HandleFunc("/api/McqEvalution", mcq.McqEvalution).Methods("GET")
	router.HandleFunc("/api/MyMcq", mcq.Mymcqquestions).Methods("POST")
	router.HandleFunc("/api/rubrics/getAll", rubrics.GetRubrics).Methods("GET")
	router.HandleFunc("/api/uploadImage", image.Upload).Methods("POST")
	router.HandleFunc("/api/serveImage/{filename}", image.ServeImage).Methods("GET")
	router.HandleFunc("/api/McqQuestions", mcq.McqQuestions).Methods("POST")

	c := cors.AllowAll()

	fmt.Print("Running....")
	handler := c.Handler(router)
	http.Handle("/", handlers.LoggingHandler(os.Stdout, handler))

	http.ListenAndServe(":8080", nil)

}
