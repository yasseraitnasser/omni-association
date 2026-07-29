package projects

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type CreateTransactionSchema struct {
	Type               string  `json:"type" validate:"required,oneof=income expense"`
	Source             string  `json:"source" validate:"required,oneof=member donor_individual government external_association"`
	DonorMemberID      *int    `json:"donor_member_id" validate:"required_if=Source member"`
	ExternalEntityName *string `json:"external_entity_name" validate:"required_if=Source external_association,required_if=Source government"`
	Amount             int     `json:"amount" validate:"gt=0"`
	PaymentMethod      string  `json:"payment_method" validate:"required,oneof=bank_transfer check cash"`
	Description        string  `json:"description" validate:"required"`
	TransactionDate    string  `json:"transaction_date" validate:"required"`
}

func validateTransactionCreationSchema(req CreateTransactionSchema) error {
	return utils.Validate.Struct(req)
}

func saveHashedFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hasher := sha256.New()
	fmt.Fprintf(hasher, "%d", time.Now().UnixNano())
	if _, err := io.Copy(hasher, src); err != nil {
		return "", err
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	hashPrefix := hex.EncodeToString(hasher.Sum(nil))[:12]
	cleanFilename := filepath.Base(fileHeader.Filename)
	newFilename := fmt.Sprintf("%s_%s", hashPrefix, cleanFilename)
	destPath := filepath.Join(utils.UPLOAD_DIR, newFilename)

	dest, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return "", err
	}

	return newFilename, nil
}

func uploadDocs(w http.ResponseWriter, r *http.Request) (string, string, error) {
	if err := os.MkdirAll(utils.UPLOAD_DIR, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return "", "", err
	}

	files := []*multipart.FileHeader{}
	filesNames := []string{"proof_doc", "receipt"}
	for _, value := range filesNames {
		_, fileHeader, err := r.FormFile(value)
		if err != nil {
			http.Error(w, "Missing or invalid file for value '"+value+"'", http.StatusBadRequest)
			return "", "", err
		}
		files = append(files, fileHeader)
	}

	savedFiles := make([]string, 0, len(files))
	for _, fileHeader := range files {
		savedName, err := saveHashedFile(fileHeader)
		if err != nil {
			http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
			return "", "", err
		}
		savedFiles = append(savedFiles, savedName)
	}

	return savedFiles[0], savedFiles[1], nil
}

func saveTransactionToDB(projectID, committeeID int, proofDocPath, receiptPath string, req CreateTransactionSchema, transactionDate time.Time) error {
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
		proofDocPath,
		receiptPath,
		req.Description,
		transactionDate,
	)
	return err
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	claims := auth.AuthenticateToken(w, r)
	if claims == nil {
		// auth.Authenticate already responses according to the error type
		return
	}

	maxSize := utils.MAX_FILE_SIZE << 20
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	metadataStr := r.FormValue("metadata")
	if metadataStr == "" {
		http.Error(w, "Missing metadata field", http.StatusBadRequest)
		return
	}

	var req CreateTransactionSchema
	if err := json.Unmarshal([]byte(metadataStr), &req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := validateTransactionCreationSchema(req)
	if err != nil {
		http.Error(w, "Invalid schema", http.StatusBadRequest)
		return
	}
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		http.Error(w, "Invalid date format: Expected 'YYYY-MM-DD'", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting project id: %v", err)
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

	proofDocPath, receiptPath, err := uploadDocs(w, r)
	if err != nil {
		return
	}

	err = saveTransactionToDB(projectID, claims.ID, proofDocPath, receiptPath, req, transactionDate)
	if err != nil {
		log.Printf("Could not insert transaction into db: %v", err)

		os.Remove(filepath.Join(utils.UPLOAD_DIR, proofDocPath))
		os.Remove(filepath.Join(utils.UPLOAD_DIR, receiptPath))

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
