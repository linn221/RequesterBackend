package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type VulnHandler struct {
	VulnService *services.VulnService
	TagService  *services.TagService
}

// Create handles POST /vulns
func (h *VulnHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[VulnInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	vuln := input.ToModel()
	vulnService, close, commit := h.VulnService.NewInstance(r.Context())
	defer close()

	id, err := vulnService.Create(r.Context(), vuln)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Connect tags if provided
	if len(input.TagIds) > 0 {
		tagService := h.TagService.CloneWithDb(vulnService.DB)
		err = tagService.ConnectTagsToReference(r.Context(), input.TagIds, "vulns", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	err = commit()
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkCreated(w, id)
}

// List handles GET /vulns
func (h *VulnHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse optional parent_id filter
	var parentId *int
	if parentIdStr := r.URL.Query().Get("parent_id"); parentIdStr != "" {
		id, err := strconv.Atoi(parentIdStr)
		if err != nil {
			utils.RespondError(w, utils.BadRequest("invalid parent_id"))
			return
		}
		parentId = &id
	}

	vulns, err := h.VulnService.List(r.Context(), parentId)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*VulnList, len(vulns))
	for i, vuln := range vulns {
		response[i] = ToVulnList(vuln)
	}

	utils.OkJson(w, response)
}

// Get handles GET /vulns/{id}
func (h *VulnHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	vuln, err := h.VulnService.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToVulnDetail(vuln))
}

// GetBySlug handles GET /vulns/slug/{slug}
func (h *VulnHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	path := r.URL.Path
	slug := path[len("/vulns/slug/"):]

	vuln, err := h.VulnService.GetBySlug(r.Context(), slug)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToVulnDetail(vuln))
}

// Update handles PUT /vulns/{id}
func (h *VulnHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	input, err := parseJson[VulnInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	vuln := input.ToModel()

	vulService, close, commit := h.VulnService.NewInstance(r.Context())
	defer close()
	_, err = vulService.Update(r.Context(), id, vuln)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	tagService := h.TagService.CloneWithDb(vulService.DB)
	// Handle tag connections for updates
	// First, get existing tags for this vuln
	existingTags, err := tagService.GetTagsForReference(r.Context(), "vulns", id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Disconnect existing tags
	for _, tag := range existingTags {
		err = tagService.DisconnectTagFromReference(r.Context(), tag.Id, "vulns", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	// Connect new tags if provided
	if len(input.TagIds) > 0 {
		err = tagService.ConnectTagsToReference(r.Context(), input.TagIds, "vulns", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}
	err = commit()
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkUpdated(w)
}

// Delete handles DELETE /vulns/{id}
func (h *VulnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	_, err = h.VulnService.Delete(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkDeleted(w)
}
