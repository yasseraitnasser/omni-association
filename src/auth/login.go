package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

func ValidateLoginSchema(req LoginSchema) error {
	validate := validator.New()
	return validate.Struct(req)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginSchema
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = ValidateLoginSchema(req)
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	match := utils.CheckPassword(req.Password, dbPassword)
	if match == false {
		log.Print("Auth fail (Password mismatch)")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := generateAccessToken(id, name, role)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
