package members

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type GetHistoryResponse struct {
	Id               int       `json:"id"`
	Project_id       int       `json:"project_id"`
	Recorded_by      int       `json:"recorded_by"`
	Type_            string    `json:"type"`
	Amount           int       `json:"amount"`
	Payment_method   string    `json:"payment_method"`
	Description      string    `json:"description"`
	Transaction_date time.Time `json:"transaction_date"`
	Created_at       time.Time `json:"created_at"`
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	claims := auth.AuthenticateToken(w, r)
	if claims == nil {
		// auth.AuthenticateToken responds with the appropriate status code
		return
	}

	query := `SELECT
	id, project_id, recorded_by,
	type, amount, payment_method,
	description, transaction_date, created_at
	FROM transactions
	WHERE donor_member_id = $1`
	rows, err := database.DB.Query(query, claims.ID)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Couldn't query the table", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var history []GetHistoryResponse
	for rows.Next() {
		var tmp GetHistoryResponse
		err := rows.Scan(
			&tmp.Id, &tmp.Project_id, &tmp.Recorded_by,
			&tmp.Type_, &tmp.Amount, &tmp.Payment_method,
			&tmp.Description, &tmp.Transaction_date, &tmp.Created_at,
		)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Couldn't scan the queried row", http.StatusInternalServerError)
			return
		}
		history = append(history, tmp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error during iteration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
