package projects

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type TransactionResponse struct {
}

type ProjectReportResponse struct {
	ProjectID                 int     `json:"project_id"`
	ProjectName               string  `json:"project_name"`
	Status                    string  `json:"status"`
	Budget                    int     `json:"budget"`
	TotalIncome               int     `json:"total_income"`
	TotalExpenses             int     `json:"total_expenses"`
	RemainingBalance          int     `json:"remaining_balance"`
	FundingProgressPercentage float64 `json:"funding_progress_percentage"`
	// Transactions              []TransactionResponse `json:"transactions"`
}

func GetReport(w http.ResponseWriter, r *http.Request) {
	claims := auth.AuthenticateToken(w, r)
	if claims == nil {
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v\n", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	query := `SELECT p.id, p.name, p.description, p.budget, p.status
	COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) AS total_income,
	COALESCE(SUM(CASE WHEN t.type = expense' THEN t.amount ELSE 0 END), 0) AS total_expenses,
	(COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) - 
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0)) AS remaining_balance
	FROM projects p
	LEFT JOIN transactions t ON p.id = t.project_id
	WHERE id = $1
	GROUP BY p.id`
	var repsone ProjectReportResponse
	err = database.DB.QueryRow(query, projectID).Scan(&repsone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such project", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// TODO: calculate the funding progress percentage

	w.WriteHeader(http.StatusOK)
}
