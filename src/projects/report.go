package projects

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type ProjectReportResponse struct {
	ProjectID                 int     `json:"project_id"`
	ProjectName               string  `json:"project_name"`
	Description               string  `json:"description"`
	Status                    string  `json:"status"`
	Budget                    int     `json:"budget"`
	TotalIncome               int     `json:"total_income"`
	TotalExpenses             int     `json:"total_expenses"`
	RemainingBalance          int     `json:"remaining_balance"`
	FundingProgressPercentage float64 `json:"funding_progress_percentage"`
}

func GetReport(w http.ResponseWriter, r *http.Request) {
	claims := auth.AuthenticateToken(w, r)
	if claims == nil {
		// auth.AuthenticateToken responds with the appropriate status code
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v\n", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	query := `SELECT p.id, p.name, p.description, p.budget, p.status,
	COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) AS total_income,
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) AS total_expenses,
	(COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) - 
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0)) AS remaining_balance
	FROM projects p
	LEFT JOIN transactions t ON p.id = t.project_id
	WHERE p.id = $1
	GROUP BY p.id`
	var response ProjectReportResponse

	err = database.DB.QueryRow(query, projectID).Scan(
		&response.ProjectID,
		&response.ProjectName,
		&response.Description,
		&response.Budget,
		&response.Status,
		&response.TotalIncome,
		&response.TotalExpenses,
		&response.RemainingBalance,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such project", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if response.Budget == 0 {
		response.FundingProgressPercentage = 100
	} else {
		response.FundingProgressPercentage = float64(response.TotalIncome) / float64(response.Budget) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
