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
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"project_name"`
	Status      string `json:"status"`
	Budget      int    `json:"budget"`
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

	query := `SELECT * FROM projects WHERE id = $1`
	var holder ProjectReportResponse
	err = database.DB.QueryRow(query, projectID).Scan(&holder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such project", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
