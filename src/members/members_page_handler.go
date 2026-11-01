package members

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type MembersPageData struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	InviteURL           string
	ExpiresAt           time.Time
	Members             []utils.Member
}

func MembersPageHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/members.html",
	)
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, "Error during parsing HTML files", http.StatusInternalServerError)
		return
	}

	query := `SELECT id, name, email, role, created_at FROM members`
	rows, err := database.DB.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No members found in the association directory")
			http.Error(w, "No members found", http.StatusInternalServerError)
			return
		}
		log.Printf("Database error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var members []utils.Member
	for rows.Next() {
		var m utils.Member
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.Role, &m.CreatedAt)
		if err != nil {
			log.Printf("Error during scanning rows: %v\n", err)
			http.Error(w, "Error during scanning rows", http.StatusInternalServerError)
			return
		}

		members = append(members, m)
	}

	data := MembersPageData{
		PageTitle:           "Omni-Association",
		IsAuthenticated:     true,
		AuthenticatedUserID: claims.ID,
		Members:             members,
	}
	tmpl.Execute(w, data)
}
