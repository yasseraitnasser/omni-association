package database

import (
	"log"

	"github.com/yasseraitnasser/omni-association/src/utils"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func AddAdminUser() {
	var err error
	adminName := utils.ADMIN_NAME
	adminEmail := utils.ADMIN_EMAIL
	adminPass := utils.ADMIN_PASS
	hash, err := HashPassword(adminPass)
	adminRole := "president"
	paidFee := true

	query := `INSERT INTO members (name, email, password, role, paid_fee)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO NOTHING;
	`
	_, err = DB.Exec(query, adminName, adminEmail, hash, adminRole, paidFee)
	if err != nil {
		log.Printf("Could not exec query: %v", err)
		log.Printf("No admin user")
		return
	}
	log.Printf("Admin added successfully: %s\n", adminEmail)
}
