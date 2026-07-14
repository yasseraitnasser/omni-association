package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

type Claims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func generateAccessToken(id int, name string, role string) (string, error) {
	duration, err := time.ParseDuration(utils.JWT_EXPIRY)
	if err != nil {
		return "", fmt.Errorf("Invalid JWT_EXPIRY format: %v", err)
	}
	claims := Claims{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.JWT_SECRET))
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

func AuthenticateToken(w http.ResponseWriter, r *http.Request) *Claims {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Unauthorized: Malformed authorization header", http.StatusUnauthorized)
		return nil
	}
	tokenString := parts[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(utils.JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
		return nil
	}
	return claims
}

func IsBoardMember(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := AuthenticateToken(w, r)
		if claims == nil {
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

func InviteMember(w http.ResponseWriter, r *http.Request) {
}
