package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yasseraitnasser/omni-association/src/database"
	"golang.org/x/crypto/bcrypt"
)

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

type CustomClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func generateAccessToken(id int, name string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiry := os.Getenv("JWT_EXPIRY")
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return "", fmt.Errorf("Invalid JWT_EXPIRY format: %v", err)
	}
	claims := CustomClaims{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateLoginSchema(req LoginSchema) error {
	validate := validator.New()
	return validate.Struct(req)
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginSchema
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("%v:", err)
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

	match := CheckPassword(req.Password, dbPassword)
	if match == false {
		log.Print("Auth fail (Password mismatch)")
		log.Print("real password: ", os.Getenv("ADMIN_PASS"))
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
