package auth

import (
	"github.com/yasseraitnasser/omni-association/src/utils"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     utils.TOKEN_NAME,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // WARN: requires HTTPS, temporarily set to false for local dev
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
