package handlers

import (
	"net/http"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type EndpointHandler struct {
	Service    *services.EndpointService
	TagService *services.TagService
}

func (h *EndpointHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[EndpointInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	id, err := h.Service.Create(r.Context(), input.ToModel())
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Connect tags if provided
	if len(input.TagIds) > 0 {
		err = h.TagService.ConnectTagsToReference(r.Context(), input.TagIds, "endpoints", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	utils.OkCreated(w, id)
}

func (h *EndpointHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	input, err := parseJson[EndpointInput](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	_, err = h.Service.Update(r.Context(), id, input.ToModel())
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Handle tag connections for updates
	// First, get existing tags for this endpoint
	existingTags, err := h.TagService.GetTagsForReference(r.Context(), "endpoints", id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Disconnect existing tags
	for _, tag := range existingTags {
		err = h.TagService.DisconnectTagFromReference(r.Context(), tag.Id, "endpoints", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	// Connect new tags if provided
	if len(input.TagIds) > 0 {
		err = h.TagService.ConnectTagsToReference(r.Context(), input.TagIds, "endpoints", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	utils.OkUpdated(w)
}

func (h *EndpointHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *EndpointHandler) List(w http.ResponseWriter, r *http.Request) {
	endpoints, err := h.Service.List(r.Context())
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	response := make([]*EndpointList, len(endpoints))
	for i, e := range endpoints {
		response[i] = ToEndpointList(e)
	}

	utils.OkJson(w, response)
}

func (h *EndpointHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	endpoint, err := h.Service.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToEndpointDetail(endpoint))
}
