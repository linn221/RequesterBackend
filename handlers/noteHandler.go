package handlers

import (
	"net/http"

	"github.com/linn221/RequesterBackend/models"
	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type NoteHandler struct {
	Service    *services.NoteService
	TagService *services.TagService
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := parseJson[CreateNoteRequest](r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	note := &models.Note{
		ReferenceType: input.ReferenceType,
		ReferenceID:   input.ReferenceId,
		Value:         input.Value,
	}

	id, err := h.Service.Create(r.Context(), note)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Connect tags if provided
	if len(input.TagIds) > 0 {
		err = h.TagService.ConnectTagsToReference(r.Context(), input.TagIds, "notes", id)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
	}

	utils.OkCreated(w, id)
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	referenceType := r.URL.Query().Get("type")
	search := r.URL.Query().Get("search")

	notes, err := h.Service.List(r.Context(), referenceType, search)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*Note, len(notes))
	for i, n := range notes {
		response[i] = ToNote(n)
	}

	utils.OkJson(w, response)
}

func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	note, err := h.Service.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToNoteDetail(note))
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	value := r.URL.Query().Get("value")
	if value == "" {
		utils.RespondError(w, utils.BadRequest("value parameter is required"))
		return
	}

	_, err = h.Service.Update(r.Context(), id, value)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkUpdated(w)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
