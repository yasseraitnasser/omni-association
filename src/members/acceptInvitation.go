package members

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type AcceptInvitationPageData struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	Token               string
	Password            string
	Error               error
}

type AcceptInvitationSchema struct {
	Token    string
	Password string
}

func validateAcceptInvitationSchema(req AcceptInvitationSchema) error {
	return utils.Validate.Struct(req)
}

func AcceptInvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing invitation token", http.StatusBadRequest)
		return
	}

	var expiresAt time.Time
	query := `SELECT invite_expiry FROM members WHERE invite_token = $1`
	err := database.DB.QueryRow(query, token).Scan(&expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid or expired invitationtoken", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if expiresAt.Before(time.Now()) {
		http.Error(w, "The invitation link has expired. Please request a new one from a board member", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/accept_invite.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := AcceptInvitationPageData{
		PageTitle:       "Omni-Association",
		IsAuthenticated: false,
		Token:           token,
	}

	tmpl.Execute(w, data)
}

func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	req := AcceptInvitationSchema{
		Token:    r.FormValue("token"),
		Password: r.FormValue("password"),
	}

	err := validateAcceptInvitationSchema(req)
	if err != nil {
		http.Error(w, "Invalid schema", http.StatusUnprocessableEntity)
		return
	}

	var expiresAt time.Time
	query := `SELECT invite_expiry FROM members WHERE invite_token = $1`
	err = database.DB.QueryRow(query, req.Token).Scan(&expiresAt)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	if expiresAt.Before(time.Now()) {
		http.Error(w, "Invalid token (token expired)", http.StatusUnauthorized)
		return
	}

	query = `UPDATE members SET password = $1, invite_token = $2, invite_expiry = $3 WHERE invite_token = $4`
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	_, err = database.DB.Exec(query, hash, nil, nil, req.Token)
	if err != nil {
		http.Error(w, "Failed to update member data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
