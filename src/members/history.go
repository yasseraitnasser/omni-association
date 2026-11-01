package members

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type History struct {
	Id               int
	Project_id       int
	Recorded_by      int
	Type_            string
	Amount           int
	Payment_method   string
	Description      string
	Transaction_date time.Time
	Created_at       time.Time
}

type profileData struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	Member              utils.Member
	History             []History
}

func GetMemberHistory(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		// auth.AuthenticateToken responds with the appropriate status code
		return
	}

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/profile.html",
	)
	if err != nil {
		log.Fatal("Error parsing html files")
	}

	vars := mux.Vars(r)
	requestedMemberID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	query := `SELECT
	id, project_id, recorded_by,
	type, amount, payment_method,
	description, transaction_date, created_at
	FROM transactions
	WHERE donor_member_id = $1`
	rows, err := database.DB.Query(query, requestedMemberID)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Couldn't query the table", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var history []History
	for rows.Next() {
		var tmp History
		err := rows.Scan(
			&tmp.Id, &tmp.Project_id, &tmp.Recorded_by,
			&tmp.Type_, &tmp.Amount, &tmp.Payment_method,
			&tmp.Description, &tmp.Transaction_date, &tmp.Created_at,
		)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Couldn't scan the queried row", http.StatusInternalServerError)
			return
		}

		history = append(history, tmp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error during iteration", http.StatusInternalServerError)
		return
	}

	query = `SELECT id, name, email, role, created_at FROM members WHERE id = $1`
	var member utils.Member

	err = database.DB.QueryRow(query, requestedMemberID).Scan(
		&member.ID,
		&member.Name,
		&member.Email,
		&member.Role,
		&member.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Not a committee member for this project", http.StatusForbidden)
			return
		}
		log.Printf("%v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	query = `SELECT amount_paid FROM membership_payments WHERE member_id = $1 AND year = $2`
	var tmp int

	member.MembershipStatus = "Active (2026 Paid)"
	err = database.DB.QueryRow(query, member.ID, 2026).Scan(&tmp)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("%v\n", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		member.MembershipStatus = "Incative"
	}

	data := profileData{
		PageTitle:           "Omni-Association",
		IsAuthenticated:     true,
		AuthenticatedUserID: claims.ID,
		Member:              member,
		History:             history,
	}
	tmpl.Execute(w, data)
}
