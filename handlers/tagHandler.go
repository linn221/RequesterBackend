package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type TagHandler struct {
	Service *services.TagService
}

// Create handles POST /tags
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[TagInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	tag := input.ToModel()
	id, err := h.Service.Create(r.Context(), tag)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkCreated(w, id)
}

// List handles GET /tags
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.Service.List(r.Context())
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*TagDTO, len(tags))
	for i, tag := range tags {
		response[i] = ToTagDTO(tag)
	}

	utils.OkJson(w, response)
}

// Get handles GET /tags/{id}
func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	tag, err := h.Service.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToTagDTO(tag))
}

// Update handles PUT /tags/{id}
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	input, err := parseJson[TagInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	tag := input.ToModel()
	_, err = h.Service.Update(r.Context(), id, tag)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkUpdated(w)
}

// ApplyTag handles POST /apply-tags/{tagId}/{referenceType}/{referenceId}
func (h *TagHandler) ApplyTag(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from URL path using PathValue
	tagIdStr := r.PathValue("tagId")
	referenceType := r.PathValue("referenceType")
	referenceIdStr := r.PathValue("referenceId")

	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		utils.RespondError(w, utils.BadRequest("invalid tag ID"))
		return
	}

	referenceId, err := strconv.Atoi(referenceIdStr)
	if err != nil {
		utils.RespondError(w, utils.BadRequest("invalid reference ID"))
		return
	}

	// Validate reference type
	validTypes := map[string]bool{
		"programs":  true,
		"endpoints": true,
		"requests":  true,
		"vulns":     true,
		"notes":     true,
	}
	if !validTypes[referenceType] {
		utils.RespondError(w, utils.BadRequest("invalid reference type"))
		return
	}

	err = h.Service.ConnectTagToReference(r.Context(), tagId, referenceType, referenceId)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkUpdated(w)
}
