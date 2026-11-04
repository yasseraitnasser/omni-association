package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

func generateAccessToken(id int, name string, role string) (string, error) {
	duration, err := time.ParseDuration(utils.JWT_EXPIRY)
	if err != nil {
		return "", fmt.Errorf("Invalid JWT_EXPIRY format: %v", err)
	}
	claims := utils.Claims{
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

func AuthenticateToken(w http.ResponseWriter, r *http.Request) (*utils.Claims, error) {
	cookie, err := r.Cookie(utils.TOKEN_NAME)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, err
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	tokenString := cookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(utils.JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		log.Print("Unauthorized: Invalid or expired token")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}
	return claims, nil
}
