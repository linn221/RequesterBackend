package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type VulnHandler struct {
	Service *services.VulnService
}

// Create handles POST /vulns
func (h *VulnHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[VulnInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	vuln := input.ToModel()
	id, err := h.Service.Create(r.Context(), vuln)
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

	vulns, err := h.Service.List(r.Context(), parentId)
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

	vuln, err := h.Service.Get(r.Context(), id)
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

	vuln, err := h.Service.GetBySlug(r.Context(), slug)
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
	_, err = h.Service.Update(r.Context(), id, vuln)
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

	_, err = h.Service.Delete(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkDeleted(w)
}
