package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type RequestHandler struct {
	Service *services.RequestService
}

func (h *RequestHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	var programId, endpointId, jobId *int
	var err error

	if programIdStr := r.URL.Query().Get("program_id"); programIdStr != "" {
		id, err := strconv.Atoi(programIdStr)
		if err != nil {
			utils.RespondError(w, utils.BadRequest("invalid program_id"))
			return
		}
		programId = &id
	}

	if endpointIdStr := r.URL.Query().Get("endpoint_id"); endpointIdStr != "" {
		id, err := strconv.Atoi(endpointIdStr)
		if err != nil {
			utils.RespondError(w, utils.BadRequest("invalid endpoint_id"))
			return
		}
		endpointId = &id
	}

	if jobIdStr := r.URL.Query().Get("job_id"); jobIdStr != "" {
		id, err := strconv.Atoi(jobIdStr)
		if err != nil {
			utils.RespondError(w, utils.BadRequest("invalid job_id"))
			return
		}
		jobId = &id
	}

	rawSQL := r.URL.Query().Get("raw_sql")
	orderBy := r.URL.Query().Get("order_by")
	asc := r.URL.Query().Get("asc") != "false" // default to true

	// Check if search query is provided
	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		requests, searchResults, err := h.Service.SearchRequests(r.Context(), searchQuery)
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		response := make([]*RequestList, len(requests))
		for i, req := range requests {
			response[i] = ToRequestList(req, searchResults)
		}

		utils.OkJson(w, response)
		return
	}

	requests, err := h.Service.List(r.Context(), programId, endpointId, jobId, rawSQL, orderBy, asc)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*RequestList, len(requests))
	for i, req := range requests {
		response[i] = ToRequestList(req, nil)
	}

	utils.OkJson(w, response)
}

func (h *RequestHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdParam(r)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	request, err := h.Service.Get(r.Context(), id)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.OkJson(w, ToRequestDetail(request))
}
