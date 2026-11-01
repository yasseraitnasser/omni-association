package members

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/database"
)

func MembershipPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: finsih implementation
	query := `INSERT INTO membership_payments (member_id, amount_paid, year, receipt_url, paid_at) VALUES ($1, $2, $3, $4, $5)`
	memberID := 1
	amountPaid := 500
	year := 2026
	receiptURL := `0fbfd28dc5da_landscape-forest-photography-trees-wallpaper-preview.jpg`
	paid_at := time.Now()

	result, err := database.DB.Exec(query, memberID, amountPaid, year, receiptURL, paid_at)
	if err != nil {
		log.Printf("Error inerting membership_payments %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	insertedRowID, err := result.LastInsertId()
	fmt.Fprintf(w, "Successfully inserted %d", insertedRowID)
}
