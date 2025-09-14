package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
)

type AttachmentHandler struct {
	AttachmentService *services.AttachmentService
}

// UploadAttachment handles POST /attachments
func (h *AttachmentHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 32MB max memory
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	// Get reference_type and reference_id from form
	referenceType := r.FormValue("reference_type")
	if referenceType == "" {
		http.Error(w, "reference_type is required", http.StatusBadRequest)
		return
	}

	referenceIDStr := r.FormValue("reference_id")
	if referenceIDStr == "" {
		http.Error(w, "reference_id is required", http.StatusBadRequest)
		return
	}

	referenceId, err := strconv.Atoi(referenceIDStr)
	if err != nil {
		http.Error(w, "invalid reference_id", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload attachment
	attachment, err := h.AttachmentService.UploadAttachment(fileHeader, referenceType, referenceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return attachment Id as plain text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%d", attachment.Id)
}

// DeleteAttachment handles DELETE /attachments?id={id}
func (h *AttachmentHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	// Get attachment Id from query parameter
	attachmentIdStr := r.URL.Query().Get("id")
	if attachmentIdStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	attachmentId, err := strconv.Atoi(attachmentIdStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	// Delete attachment
	if err := h.AttachmentService.DeleteAttachment(attachmentId); err != nil {
		if err.Error() == "attachment not found" {
			http.Error(w, "attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// GetAttachment handles GET /attachments/{id}
func (h *AttachmentHandler) GetAttachment(w http.ResponseWriter, r *http.Request) {
	// Extract Id from URL path
	path := r.URL.Path
	idStr := path[len("/attachments/"):]

	attachmentId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid attachment id", http.StatusBadRequest)
		return
	}

	// Get attachment
	attachment, err := h.AttachmentService.GetAttachment(attachmentId)
	if err != nil {
		if err.Error() == "attachment not found" {
			http.Error(w, "attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := Attachment{
		Id:       int(attachment.Id),
		Filename: attachment.Filename,
		URL:      attachment.GetURL(),
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ServeAttachment serves the actual file content
func (h *AttachmentHandler) ServeAttachment(w http.ResponseWriter, r *http.Request) {
	// Extract Id from URL path
	path := r.URL.Path
	idStr := path[len("/attachments/"):]

	attachmentId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid attachment id", http.StatusBadRequest)
		return
	}

	// Get attachment
	attachment, err := h.AttachmentService.GetAttachment(attachmentId)
	if err != nil {
		if err.Error() == "attachment not found" {
			http.Error(w, "attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve file
	if err := h.AttachmentService.ServeFile(w, r, attachment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
