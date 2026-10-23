package auth

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

func ValidateLoginSchema(req LoginSchema) error {
	return utils.Validate.Struct(req)
}

type loginPageData struct {
	PageTitle string
	Error     string
	Email     string
}

func renderLoginError(w http.ResponseWriter, r *http.Request, email string) {
	w.WriteHeader(http.StatusUnauthorized)

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/login.html",
	)
	if err != nil {
		log.Fatal("Error parsing html files")
	}

	data := loginPageData{
		PageTitle: "Omni-Association",
		Error:     "Invalid email or password.",
		Email:     email,
	}
	log.Print(data.PageTitle)
	log.Print(data.Email)
	log.Print(data.Error)
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req := LoginSchema{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	err := ValidateLoginSchema(req)
	if err != nil {
		http.Error(w, "Invalid Schema", http.StatusUnprocessableEntity)
		return
	}

	var (
		id         int
		name       string
		dbPassword string
		role       string
	)

	query := `SELECT id, name, password, role FROM members WHERE email = $1`
	err = database.DB.QueryRow(query, req.Email).Scan(&id, &name, &dbPassword, &role)
	if err != nil {
		log.Printf("Auth fail (User not found): %v", err)
		renderLoginError(w, r, req.Email)
		return
	}

	match := utils.CheckPassword(req.Password, dbPassword)
	if match == false {
		log.Print("Auth fail (Password mismatch)")
		renderLoginError(w, r, req.Email)
		return
	}

	token, err := generateAccessToken(id, name, role)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	maxAgeInHours, err := time.ParseDuration(utils.JWT_EXPIRY)
	if err != nil {
		http.Error(w, "Internal Server ERROR", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     utils.TOKEN_NAME,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // WARN: requires HTTPS, temporarily set to false for local dev
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(maxAgeInHours.Hours()) * 3600,
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/login.html",
	)
	if err != nil {
		log.Fatal("Error parsing html files")
	}

	data := loginPageData{
		PageTitle: "Omni-Association",
	}
	tmpl.Execute(w, data)
}
