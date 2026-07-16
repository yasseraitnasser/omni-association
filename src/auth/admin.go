package auth

import (
	"log"

	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

func AddAdminUser() {
	var err error
	adminName := utils.ADMIN_NAME
	adminEmail := utils.ADMIN_EMAIL
	adminPass := utils.ADMIN_PASS
	hash, err := HashPassword(adminPass)
	if err != nil {
		log.Printf("Could not hash password")
		return
	}
	adminRole := "president"
	paidFee := true

	query := `INSERT INTO members (name, email, password, role, paid_fee)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO NOTHING;
	`
	_, err = database.DB.Exec(query, adminName, adminEmail, hash, adminRole, paidFee)
	if err != nil {
		log.Printf("Could not exec query: %v", err)
		log.Printf("No admin user")
		return
	}
	log.Printf("Admin added successfully: %s\n", adminEmail)
}
