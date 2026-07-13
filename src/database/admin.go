package database

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AddAdminUser() {
	var err error
	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASS")
	hash, err := hashPassword(adminPass)
	adminRole := "president"
	paidFee := true

	query := `
		INSERT INTO members
		(name, email, password, role, paid_fee)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO NOTHING;
	`
	_, err = DB.Exec(query, adminName, adminEmail, hash, adminRole, paidFee)
	if err != nil {
		log.Printf("Could not exec query: %v", err)
	}
	log.Printf("Admin added successfully: %s\n", adminEmail)
}
