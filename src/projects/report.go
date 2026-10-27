package projects

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type Project struct {
	ID                        int
	Name                      string
	Description               string
	Status                    string
	Budget                    int
	TotalIncome               int
	TotalExpenses             int
	RemainingBalance          int
	FundingProgressPercentage float64
	CreatedAt                 time.Time
}

type ProjectsPageData struct {
	PageTitle string
	Projects  []Project
}

func GetReportReport(w http.ResponseWriter, r *http.Request) {
	_, err := auth.AuthenticateToken(w, r)
	if err != nil {
		// auth.AuthenticateToken responds with the appropriate status code
		return
	}

	vars := mux.Vars(r)
	_, err = strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v\n", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}
	// TODO: finsih handler implementation
}

func GetProjects() ([]Project, error) {
	query := `SELECT p.id, p.name, p.description, p.budget, p.status, p.created_at,
	COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) AS total_income,
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) AS total_expenses,
	(COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) - 
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0)) AS remaining_balance
	FROM projects p
	LEFT JOIN transactions t ON p.id = t.project_id
	GROUP BY p.id
	ORDER BY created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var response Project
	// if response.Budget == 0 {
	// 	response.FundingProgressPercentage = 100
	// } else {
	// 	response.FundingProgressPercentage = float64(response.TotalIncome) / float64(response.Budget) * 100
	// }

	var projectsList []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name,
			&p.Description, &p.Budget,
			&p.Status, &p.CreatedAt,
			&p.TotalIncome, &p.TotalExpenses,
			&p.RemainingBalance); err != nil {
			return nil, err
		}
		projectsList = append(projectsList, p)
	}

	return projectsList, nil
}
