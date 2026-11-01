package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

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
		log.Printf("%v\n", err)
		http.Error(w, "Error during parsing HTML files", http.StatusInternalServerError)
		return
	}

	data := utils.Data{
		PageTitle:       "Omni-Association",
		IsAuthenticated: false,
	}
	tmpl.Execute(w, data)
}

func getMemberData(w http.ResponseWriter, claims *utils.Claims) (utils.Member, error) {
	query := `SELECT name, role, created_at FROM members WHERE id = $1`
	var member utils.Member

	err := database.DB.QueryRow(query, claims.ID).Scan(&member.Name, &member.Role, &member.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such member", http.StatusForbidden)
			return member, err
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return member, err
	}
	return member, nil
}

type MembershipPayment struct {
	MemberID   int
	AmountPaid int
	Year       int
	PaidAt     time.Time
}

func getMembershipPaymentsData(w http.ResponseWriter, claims *utils.Claims, member *utils.Member) error {
	query := `SELECT member_id, amount_paid, year, paid_at FROM membership_payments WHERE member_id = $1`
	rows, err := database.DB.Query(query, claims.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such member", http.StatusForbidden)
			return err
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return err
	}
	defer rows.Close()

	member.TotalPaid = 0
	member.MembershipStatus = "Inactive"
	for rows.Next() {
		var mp MembershipPayment
		if err := rows.Scan(&mp.MemberID, &mp.AmountPaid, &mp.Year, &mp.PaidAt); err != nil {
			log.Printf("Error scanning membership payments: %v\n", err)
			http.Error(w, "Failed to load membership payments", http.StatusInternalServerError)
			return err
		}

		member.TotalPaid += mp.AmountPaid

		if mp.Year == time.Now().Year() {
			member.MembershipStatus = fmt.Sprintf("Active (%d Paid)", mp.Year)
		}
	}
	return nil
}

func getAssociationFee(w http.ResponseWriter) (int, error) {
	query := `SELECT current_fee FROM association WHERE id = $1`
	var associationFee int
	err := database.DB.QueryRow(query, 3).Scan(&associationFee)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such association", http.StatusForbidden)
			return 0, err
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return 0, err
	}
	return associationFee, nil
}

func getProjectCount(w http.ResponseWriter) (int, error) {
	query := `SELECT id FROM projects`
	projectCount := 0
	rows, err := database.DB.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such association", http.StatusForbidden)
			return 0, err
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		projectCount++
	}
	return projectCount, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/dashboard.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	member, err := getMemberData(w, claims)
	if err != nil {
		return
	}

	err = getMembershipPaymentsData(w, claims, &member)
	if err != nil {
		return
	}

	associationFee, err := getAssociationFee(w)
	if err != nil {
		return
	}

	projectCount, err := getProjectCount(w)
	if err != nil {
		return
	}

	data := utils.DashboardData{
		PageTitle:           "Omni-Association",
		IsAuthenticated:     true,
		AuthenticatedUserID: claims.ID,
		AssociationFee:      associationFee,
		ProjectCount:        projectCount,
		Member:              member,
	}
	tmpl.Execute(w, data)
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

func RegisterMembershipPayment(w http.ResponseWriter, r *http.Request) {
	query := `INSERT INTO membership_payments (member_id, amount_paid, year, receipt_url, paid_at) VALUES ($1, $2, $3, $4, $5)`
	memberID := 1
	amountPaid := 500
	year := 2026
	receiptURL := `0fbfd28dc5da_landscape-forest-photography-trees-wallpaper-preview.jpg`
	paid_at := time.Now()

	result, err := database.DB.Exec(query, memberID, amountPaid, year, receiptURL, paid_at)
	if err != nil {
		log.Printf("Error inerting membership_payments %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	insertedRowID, err := result.LastInsertId()
	fmt.Fprintf(w, "Successfully inserted %d", insertedRowID)
}

func createAssociation() {
	query := `SELECT id FROM association WHERE id = $1`

	_, err := database.DB.Query(query, 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query = `INSERT INTO association (name, current_fee) VALUES ($1, $2)`
			name := "Omni-Association"
			currentFee := 0

			_, err = database.DB.Exec(query, name, currentFee)
			if err != nil {
				log.Fatalf("Error inerting %s into association: %v\n", name, err)
			}

			log.Printf("Successfully inserted %s into the database", name)
			return
		}
		log.Printf("Database error: %v\n", err)
		return
	}
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
	router.HandleFunc("/login", auth.Login).Methods("POST")

	router.HandleFunc("/api/members", members.MembersPageHandler).Methods("GET")
	router.HandleFunc("/api/members/invite", Chain(members.InviteMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/members/accept", members.AcceptInvitation).Methods("POST")

	router.HandleFunc("/projects", projects.ProjectsPageHandler).Methods("GET")
	router.HandleFunc("/api/projects", Chain(projects.CreateProject, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/projects/{id:[0-9]+}/committee", Chain(projects.AssignCommitteeMember, members.IsBoardMember)).Methods("POST")
	router.HandleFunc("/api/projects/{id:[0-9]+}/transactions", projects.CreateTransaction).Methods("POST")
	router.HandleFunc("/api/projects/{id:[0-9]+}/report", projects.ProjectDetailsHandler).Methods("GET")

	router.HandleFunc("/api/membership_payment", RegisterMembershipPayment).Methods("POST")

	router.HandleFunc("/api/members/{id:[0-9]+}/history", members.GetMemberHistory).Methods("GET")

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	createAssociation()

	auth.AddAdminUser()
	utils.InitValidator()

	domain := utils.SERVER_HOST + ":" + utils.SERVER_PORT
	log.Printf("Listening on: %s\n", domain)
	http.ListenAndServe(domain, router)
}
