package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
)

type ImageHandler struct {
	ImageService *services.ImageService
}

// UploadImages handles POST /images
func (h *ImageHandler) UploadImages(w http.ResponseWriter, r *http.Request) {
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

	// Get uploaded files (multiple files)
	form := r.MultipartForm
	files := form.File["files"]
	if len(files) == 0 {
		http.Error(w, "at least one file is required", http.StatusBadRequest)
		return
	}

	// Upload images
	images, err := h.ImageService.UploadImages(files, referenceType, referenceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response format
	response := make([]Image, len(images))
	for i, image := range images {
		response[i] = Image{
			Id:           int(image.Id),
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.GetURL(),
		}
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteImage handles DELETE /images?id={id}
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	// Get image Id from query parameter
	imageIdStr := r.URL.Query().Get("id")
	if imageIdStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	// Delete image
	if err := h.ImageService.DeleteImage(imageId); err != nil {
		if err.Error() == "image not found" {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// GetImage handles GET /images/{id}
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	// Extract Id from URL path
	path := r.URL.Path
	idStr := path[len("/images/"):]

	imageId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid image id", http.StatusBadRequest)
		return
	}

	// Get image
	image, err := h.ImageService.GetImage(imageId)
	if err != nil {
		if err.Error() == "image not found" {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := Image{
		Id:           int(image.Id),
		Filename:     image.Filename,
		OriginalName: image.OriginalName,
		URL:          image.GetURL(),
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ServeImage serves the actual image file content
func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	path := r.URL.Path
	filename := path[len("/images/file/"):]

	// Get image by filename
	image, err := h.ImageService.GetImageByFilename(filename)
	if err != nil {
		if err.Error() == "image not found" {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve file
	if err := h.ImageService.ServeImage(w, r, image); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
