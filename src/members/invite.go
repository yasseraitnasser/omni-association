package members

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type InviteMemberSchema struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name"`
	Role  string `json:"role" validate:"required,oneof=vice-president treasurer assistant-treasurer general-secretary assistant-general-secretary advisor member"`
}

func IsBoardMemberMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.AuthenticateToken(w, r)
		if err != nil {
			return
		}

		isBoardMember := claims.Role == "president" ||
			claims.Role == "vice-president" ||
			claims.Role == "treasurer" ||
			claims.Role == "assistant-treasurer" ||
			claims.Role == "general-secretary" ||
			claims.Role == "assistant-general-secretary" ||
			claims.Role == "advisor"
		if !isBoardMember {
			http.Error(w, "Forbidden: Board Member privileges required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func validateMemberInvitationSchema(req InviteMemberSchema) error {
	return utils.Validate.Struct(req)
}

func GenerateSecureToken() (string, error) {
	bytes := make([]byte, utils.SECURE_TOKEN_LENGTH)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func savePendingMemberToDB(name, email, role, token string, expiry time.Time) error {
	query := `INSERT INTO members (name, email, role, invite_token, invite_expiry)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := database.DB.Exec(query, name, email, role, token, expiry)
	if err != nil {
		log.Printf("Could not exec query: %v\n", err)
		return err
	}
	log.Printf("Member added successfully: %s\n", email)
	return nil
}

func InviteMember(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	req := InviteMemberSchema{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
		Role:  r.FormValue("role"),
	}

	err = validateMemberInvitationSchema(req)
	if err != nil {
		http.Error(w, "Invalid Schema", http.StatusUnprocessableEntity)
		return
	}

	token, err := GenerateSecureToken()
	if err != nil {
		log.Printf("Couldn't generate secure token: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	duration, err := time.ParseDuration(utils.SECURE_TOKEN_EXPIRY)
	if err != nil {
		duration = 24 * time.Hour
	}
	expiry := time.Now().Add(duration)

	err = savePendingMemberToDB(req.Name, req.Email, req.Role, token, expiry)
	if err != nil {
		http.Error(w, "Failed to create invitation: email might already exist", http.StatusConflict)
		return
	}

	inviteURL := "http://" + utils.SERVER_HOST + "/accept-invite?token=" + token

	RenderMembersPage(w, claims, inviteURL, expiry)
}
