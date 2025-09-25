package handlers

import (
	"net/http"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type ProgramHandler struct {
	Service    *services.ProgramService
	TagService *services.TagService
}

func (h *ProgramHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[ProgramInput](r)
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
		err = h.TagService.ConnectTagsToReference(r.Context(), input.TagIds, "programs", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	utils.OkCreated(w, id)
}

func (h *ProgramHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	input, err := parseJson[ProgramInput](r)
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
	// First, get existing tags for this program
	existingTags, err := h.TagService.GetTagsForReference(r.Context(), "programs", id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Disconnect existing tags
	for _, tag := range existingTags {
		err = h.TagService.DisconnectTagFromReference(r.Context(), tag.Id, "programs", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	// Connect new tags if provided
	if len(input.TagIds) > 0 {
		err = h.TagService.ConnectTagsToReference(r.Context(), input.TagIds, "programs", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	utils.OkUpdated(w)
}

func (h *ProgramHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *ProgramHandler) List(w http.ResponseWriter, r *http.Request) {
	programs, err := h.Service.List(r.Context())
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	response := make([]*ProgramList, len(programs))
	for i, p := range programs {
		response[i] = ToProgramList(p)
	}

	utils.OkJson(w, response)
}

func (h *ProgramHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}
	program, err := h.Service.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToProgramDetail(program))
}
