package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/members"
	"github.com/yasseraitnasser/omni-association/src/projects"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type Data struct {
	PageTitle string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/layouts/layout.html")
	if err != nil {
		log.Fatal("Error parsing layout.html")
	}
	data := Data{
		PageTitle: "omni-association",
	}
	tmpl.Execute(w, data)
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

func main() {
	router := mux.NewRouter()
	var err error
	err = utils.InitEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router.HandleFunc("/login", (auth.Login)).Methods("POST")
	router.HandleFunc("/api/members/invite", Chain(members.InviteMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/members/accept", members.AcceptInvitation).Methods("POST")
	router.HandleFunc("/api/projects", Chain(projects.CreateProject, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/projects/{id}/committee", Chain(projects.AssignCommitteeMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/", home)

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	auth.AddAdminUser()

	domain := utils.SERVER_HOST + ":" + utils.SERVER_PORT
	log.Printf("Listening on: %s\n", domain)
	http.ListenAndServe(domain, router)
}
