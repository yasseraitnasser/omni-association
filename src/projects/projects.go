package projects

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type CreateProjectSchema struct {
	ProjectName string `json:"project_name" validate:"required"`
	Description string `json:"description" validate:"required"`
	LeaderID    int    `json:"leader_id" validate:"required"`
	Budget      int    `json:"budget" validate:"required"`
}

func validateProjectCreationSchema(req CreateProjectSchema) error {
	validate := validator.New()
	return validate.Struct(req)
}

func saveProjectToDB(projectName, descrption string, leaderID, budget int) error {
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO projects (name, description, budget) VALUES ($1, $2, $3) RETURNING id;`
	var projectID int
	err = tx.QueryRowContext(ctx, query, projectName, descrption, budget).Scan(&projectID)
	if err != nil {
		log.Printf("Could not insert project: %v\n", err)
		return err
	}

	query = `INSERT INTO project_committees (project_id, member_id, role_in_project) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, query, projectID, leaderID, "project-lead")
	if err != nil {
		log.Printf("Could not insert project committee member: %v\n", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return err
	}

	log.Printf("Project added successfully: %s\n", projectName)
	return nil
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectSchema
	var err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = validateProjectCreationSchema(req)
	if err != nil {
		http.Error(w, "Invalid Schema", http.StatusBadRequest)
		return
	}

	query := `SELECT id FROM members WHERE id = $1`
	var holder int
	err = database.DB.QueryRow(query, req.LeaderID).Scan(&holder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No such member", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = saveProjectToDB(req.ProjectName, req.Description, req.LeaderID, req.Budget)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type AssignCommitteeSchema struct {
	Name string
	Role string
}

func AssignCommitteeMember(w http.ResponseWriter, r *http.Request) {
}
