package members

import (
	"log"
	"net/http"
	"strconv"

	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
	"github.com/yasseraitnasser/omni-association/src/projects"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

type MembershipPaymentData struct {
	MemberID   int
	AmountPaid int
	Year       int
	ReceiptURL string
}

func validateMembershipPaymentSchema(req MembershipPaymentData) error {
	return utils.Validate.Struct(req)
}

func convertPaymentData(r *http.Request) (int, int, int, error) {
	memberID, err := strconv.Atoi(r.FormValue("member_id"))
	if err != nil {
		return 0, 0, 0, err
	}

	amountPaid, err := strconv.Atoi(r.FormValue("amount_paid"))
	if err != nil {
		return 0, 0, 0, err
	}

	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		return 0, 0, 0, err
	}

	return memberID, amountPaid, year, nil
}

func MembershipPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	isBoardMember := utils.IsBoardMember(*claims)
	if !isBoardMember {
		http.Error(w, "Board Member privilege required", http.StatusUnauthorized)
		return
	}

	maxSize := utils.MAX_FILE_SIZE << 20
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	memberID, amountPaid, year, err := convertPaymentData(r)
	if err != nil {
		http.Error(w, "Enable to convert form data", http.StatusBadRequest)
		return
	}

	_, header, err := r.FormFile("receipt")
	if err != nil {
		http.Error(w, "Missing or invalid file for 'receipt'", http.StatusBadRequest)
		return
	}

	savedName, err := projects.SaveHashedFile(header)
	if err != nil {
		log.Printf("Failed to save file: %v\n", err)
		http.Error(w, "Files to save file", http.StatusInternalServerError)
		return
	}

	req := MembershipPaymentData{
		MemberID:   memberID,
		AmountPaid: amountPaid,
		Year:       year,
		ReceiptURL: savedName,
	}

	err = validateMembershipPaymentSchema(req)
	if err != nil {
		http.Error(w, "Invalid schema", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO membership_payments (member_id, amount_paid, year, receipt_url) VALUES ($1, $2, $3, $4)`

	_, err = database.DB.Exec(query, memberID, amountPaid, year, savedName)
	if err != nil {
		log.Printf("Error inserting membership_payments %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/api/members/", http.StatusCreated)
}
