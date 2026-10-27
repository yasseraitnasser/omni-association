package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/members"
	"github.com/yasseraitnasser/omni-association/src/projects"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/home.html",
	)
	if err != nil {
		log.Fatal("Error parsing html files")
	}
	data := utils.Data{
		PageTitle: "Omni-Association",
	}
	tmpl.Execute(w, data)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("omni_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString := cookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(utils.JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		log.Print("Unauthorized: Invalid or expired token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/dashboard.html",
	)
	if err != nil {
		log.Fatal("Error parsing html files")
	}

	data := utils.DashboardData{
		PageTitle:      "Omni-Association",
		AssociationFee: 500,
		Member: struct {
			Name         string
			Role         string
			CreatedAt    time.Time
			TotalPaid    int
			ProjectCount int
		}{
			Name: claims.Name,
		},
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

	staticDir := "./src/templates/layouts/pico-main/css"
	staticHandler := http.StripPrefix("/pico-main/css/", http.FileServer(http.Dir(staticDir)))
	router.PathPrefix("/pico-main/css/").Handler(staticHandler)

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")

	router.HandleFunc("/login", auth.LoginPageHandler).Methods("GET")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("GET")
	router.HandleFunc("/login", (auth.Login)).Methods("POST")

	router.HandleFunc("/api/members/invite", Chain(members.InviteMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/members/accept", members.AcceptInvitation).Methods("POST")

	router.HandleFunc("/projects", projects.ProjectsPageHandler).Methods("GET")
	router.HandleFunc("/api/projects", Chain(projects.CreateProject, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/projects/{id:[0-9]+}/committee", Chain(projects.AssignCommitteeMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/projects/{id:[0-9]+}/transactions", projects.CreateTransaction).Methods("POST")
	// router.HandleFunc("/api/projects/{id:[0-9]+}/report", projects.GetReport).Methods("GET") TODO: implement project details handler

	router.HandleFunc("/api/members/me/history", members.GetHistory).Methods("GET")

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	auth.AddAdminUser()
	utils.InitValidator()

	domain := utils.SERVER_HOST + ":" + utils.SERVER_PORT
	log.Printf("Listening on: %s\n", domain)
	http.ListenAndServe(domain, router)
}
