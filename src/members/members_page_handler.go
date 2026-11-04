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
	IsBoardMember       bool
}

func RenderMembersPage(w http.ResponseWriter, claims *utils.Claims, inviteURL string, expiresAt time.Time) {
	isBoardMember := utils.IsBoardMember(*claims)

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
		m.HasPaidCurrentYear = true

		query = `SELECT id FROM membership_payments WHERE member_id = $1 AND year = $2`
		var tmp int

		err = database.DB.QueryRow(query, m.ID, 2026).Scan(&tmp)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			m.HasPaidCurrentYear = false
		}

		members = append(members, m)
	}

	data := MembersPageData{
		PageTitle:           "Omni-Association",
		IsAuthenticated:     true,
		AuthenticatedUserID: claims.ID,
		Members:             members,
		IsBoardMember:       isBoardMember,
		InviteURL:           inviteURL,
		ExpiresAt:           expiresAt,
	}
	tmpl.Execute(w, data)
}

func MembersPageHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	RenderMembersPage(w, claims, "", time.Now())
}
