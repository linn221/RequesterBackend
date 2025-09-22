package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type ImportBurpHandler struct {
	Service *services.ImportBurpService
}

func (h *ImportBurpHandler) ImportBurpXML(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max file size
	if err != nil {
		utils.RespondError(w, utils.BadRequest("failed to parse multipart form"))
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondError(w, utils.BadRequest("file is required"))
		return
	}
	defer file.Close()

	// Get filename
	filename := handler.Filename

	// Parse additional form fields
	programIdStr := r.FormValue("program_id")
	if programIdStr == "" {
		utils.RespondError(w, utils.BadRequest("program_id is required"))
		return
	}

	programId, err := strconv.Atoi(programIdStr)
	if err != nil {
		utils.RespondError(w, utils.BadRequest("program_id must be a valid integer"))
		return
	}

	ignoredHeaders := r.FormValue("ignored_headers")

	// Import the Burp XML file
	jobId, err := h.Service.ImportBurpXML(r.Context(), file, filename, programId, ignoredHeaders)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Return the job ID as plain text
	utils.OkCreated(w, jobId)
}
