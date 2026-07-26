package projects

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type CreateTransactionSchema struct {
	Type               string  `json:"type" validate:"required,oneof=income expense"`
	Source             string  `json:"source" validate:"required,oneof=member donor_individual government external_association"`
	DonorMemberID      *int    `json:"donor_member_id" validate:"required_if=Source member"`
	ExternalEntityName *string `json:"external_entity_name" validate:"required_if=Source external_association,required_if=Source government"`
	Amount             int     `json:"amount" validate:"gt=0"`
	PaymentMethod      string  `json:"payment_method" validate:"required,oneof=bank_transfer check cash"`
	ProofDocURL        string  `json:"proof_doc_url" validate:"required,url"`
	ReceiptURL         string  `json:"receipt_url" validate:"required,url"`
	Description        string  `json:"description" validate:"required"`
	TransactionDate    string  `json:"transaction_date" validate:"required"`
}

func validateTransactionCreationSchema(req CreateTransactionSchema) error {
	validate := validator.New()
	return validate.Struct(req)
}

func saveTransactionToDB(projectID, committeeID int, req CreateTransactionSchema, transactionDate time.Time) error {
	query := `INSERT INTO transactions (
		project_id,
		recorded_by,
		type,
		source,
		donor_member_id,
		external_entity_name,
		amount,
		payment_method,
		proof_doc_url,
		receipt_url,
		description,
		transaction_date
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := database.DB.Exec(
		query,
		projectID,
		committeeID,
		req.Type,
		req.Source,
		req.DonorMemberID,
		req.ExternalEntityName,
		req.Amount,
		req.PaymentMethod,
		req.ProofDocURL,
		req.ReceiptURL,
		req.Description,
		transactionDate,
	)
	return err
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionSchema
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = validateTransactionCreationSchema(req)
	if err != nil {
		http.Error(w, "Invalid schema", http.StatusBadRequest)
		return
	}
	transactionDate, err := time.Parse("2006-01-02 at 15:04", req.TransactionDate)
	if err != nil {
		http.Error(w, "Invalid date format: Expected 'YYYY-MM-DD at HH:MM'", http.StatusBadRequest)
		return
	}

	claims := auth.AuthenticateToken(w, r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting id: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	query := `SELECT role_in_project FROM project_committees WHERE project_id = $1 AND member_id = $2`
	var roleInProject string
	err = database.DB.QueryRow(query, projectID, claims.ID).Scan(&roleInProject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Not a committee member for this project", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = saveTransactionToDB(projectID, claims.ID, req, transactionDate)
	if err != nil {
		log.Printf("Could not insert transaction into db: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
