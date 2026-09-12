package members

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type AcceptInvitationSchema struct {
	Token    string `json:"token"`
	Password string `json:"password" validate:"required,min=12"`
}

func validateAcceptInvitationSchema(req AcceptInvitationSchema) error {
	return utils.Validate.Struct(req)
}

func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req AcceptInvitationSchema
	var err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = validateAcceptInvitationSchema(req)
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

	w.WriteHeader(http.StatusOK)
}
