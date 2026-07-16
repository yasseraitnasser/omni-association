package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

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

func AuthenticateToken(w http.ResponseWriter, r *http.Request) *Claims {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return nil
	}

	parts := strings.SplitAfterN(authHeader, " ", 2) // NOTE: <=> strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer " {
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
