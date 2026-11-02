package projects

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type Committee struct {
	ID            int
	Name          string
	RoleInProject string
	JoinAt        time.Time
}

type Transaction struct {
	ID                 int
	ProjectId          int
	RecordedBy         int
	Type               string
	Source             string
	DonorMemberId      *int
	ExternalEntityName *string
	Amount             int
	PaymentMethod      string
	ProofDocUrl        string
	ReceiptUrl         string
	Description        string
	TransactionDate    time.Time
	CreatedAt          time.Time
}

type ProjectDetails struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	Project             utils.Project
	NetBalance          int
	Committee           []Committee
	Transactions        []Transaction
	IsBoardMember       bool
}

func getProjectDetails(w http.ResponseWriter, projectID int) (*utils.Project, error) {
	query := `SELECT p.id, p.name, p.description, p.budget, p.status, p.created_at,
	COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) AS total_income,
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) AS total_expenses,
	(COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) - 
	COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0)) AS remaining_balance
	FROM projects p
	LEFT JOIN transactions t ON p.id = t.project_id
	WHERE p.id = $1
	GROUP BY p.id
	ORDER BY created_at DESC`
	var project utils.Project

	err := database.DB.QueryRow(query, projectID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Budget,
		&project.Status,
		&project.CreatedAt,
		&project.TotalIncome,
		&project.TotalExpenses,
		&project.RemainingBalance,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such project", http.StatusForbidden)
			return nil, err
		}
		log.Printf("%v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return nil, err
	}
	return &project, nil
}

func getProjectCommitteeDetails(w http.ResponseWriter, projectID int) ([]Committee, error) {
	query := `SELECT member_id, role_in_project, join_at FROM project_committees WHERE project_id = $1`
	rows, err := database.DB.Query(query, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such project_committee", http.StatusForbidden)
			return nil, err
		}
		log.Printf("%v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	var committees []Committee
	for rows.Next() {
		var c Committee
		err = rows.Scan(&c.ID, &c.RoleInProject, &c.JoinAt)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Error during scanning", http.StatusInternalServerError)
			return nil, err
		}

		query = `SELECT name FROM members WHERE id = $1`
		err = database.DB.QueryRow(query, c.ID).Scan(&c.Name)

		committees = append(committees, c)
	}
	return committees, nil
}

func getProjectTransactions(w http.ResponseWriter, projectID int) ([]Transaction, error) {
	query := `SELECT id, project_id, recorded_by, type, source, donor_member_id, external_entity_name, amount, payment_method, proof_doc_url, receipt_url, description, transaction_date, created_at FROM transactions WHERE project_id = $1`

	rows, err := database.DB.Query(query, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such transaction", http.StatusForbidden)
			return nil, err
		}
		log.Printf("%v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err = rows.Scan(
			&t.ID,
			&t.ProjectId,
			&t.RecordedBy,
			&t.Type,
			&t.Source,
			&t.DonorMemberId,
			&t.ExternalEntityName,
			&t.Amount,
			&t.PaymentMethod,
			&t.ProofDocUrl,
			&t.ReceiptUrl,
			&t.Description,
			&t.TransactionDate,
			&t.CreatedAt,
		)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Error during scanning", http.StatusInternalServerError)
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func ProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		// auth.AuthenticateToken responds with the appropriate status code
		return
	}

	isBoardMember := utils.IsBoardMember(*claims)

	tmpl, err := template.ParseFiles(
		"src/templates/layouts/layout.html",
		"src/templates/pages/project_details.html",
	)
	if err != nil {
		log.Printf("Error parsing templates: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v\n", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := getProjectDetails(w, projectID)
	if err != nil {
		return
	}

	committees, err := getProjectCommitteeDetails(w, projectID)
	if err != nil {
		return
	}

	transactions, err := getProjectTransactions(w, projectID)
	if err != nil {
		return
	}

	data := ProjectDetails{
		PageTitle:           "Omni-Association",
		IsAuthenticated:     true,
		AuthenticatedUserID: claims.ID,
		Project:             *project,
		Committee:           committees,
		Transactions:        transactions,
		IsBoardMember:       isBoardMember,
	}

	tmpl.Execute(w, data)
}
