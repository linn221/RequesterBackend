package handlers

import (
	"net/http"

	"github.com/linn221/RequesterBackend/utils"
)

type Handler struct{}

// ===== Notes =====
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// ===== Attachments =====
func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// ===== Programs =====
func (h *Handler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) GetProgram(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// ===== Endpoints =====
func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) ListEndpoints(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) GetEndpoint(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// ===== Requests =====
func (h *Handler) ListRequests(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// ===== Import HAR =====
func (h *Handler) ImportHAR(w http.ResponseWriter, r *http.Request) {
	// This method is kept for backward compatibility but should not be used
	// Use ImportHarHandler instead
	utils.RespondError(w, utils.BadRequest("use dedicated ImportHarHandler"))
}

// ===== Jobs =====
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
