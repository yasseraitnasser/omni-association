package projects

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/utils"
)

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.AuthenticateToken(w, r)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	parsedFileName := filepath.Base(filename)
	filePath := filepath.Join(utils.UPLOAD_DIR, parsedFileName)

	http.ServeFile(w, r, filePath)
}
