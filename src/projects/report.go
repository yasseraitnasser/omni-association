package projects

import (
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type ProjectsPageData struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	IsBoardMember       bool
	Projects            []utils.Project
}

func GetProjects() ([]utils.Project, error) {
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
	// } // TODO: calculate funding progress percentage

	var projectsList []utils.Project
	for rows.Next() {
		var p utils.Project
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
