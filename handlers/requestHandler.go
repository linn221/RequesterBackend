package handlers

import (
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/models"
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

	// Parse new filter parameters
	domain := r.URL.Query().Get("domain")
	urlContains := r.URL.Query().Get("url_contains")
	urlMatch := r.URL.Query().Get("url_match")
	includeSubdomains := r.URL.Query().Get("includeSubdomains") == "1"

	// Parse multi-level ordering parameters
	orderBy1 := r.URL.Query().Get("order_by1")
	asc1 := r.URL.Query().Get("asc1") != "false" // default to true
	orderBy2 := r.URL.Query().Get("order_by2")
	asc2 := r.URL.Query().Get("asc2") != "false" // default to true
	orderBy3 := r.URL.Query().Get("order_by3")
	asc3 := r.URL.Query().Get("asc3") != "false" // default to true
	orderBy4 := r.URL.Query().Get("order_by4")
	asc4 := r.URL.Query().Get("asc4") != "false" // default to true

	// Check if search query is provided
	searchQuery := r.URL.Query().Get("search")
	var requests []*models.MyRequest

	if searchQuery != "" {
		requests, err = h.Service.SearchRequests(r.Context(), searchQuery, domain, urlContains, urlMatch, includeSubdomains, orderBy1, asc1, orderBy2, asc2, orderBy3, asc3, orderBy4, asc4)
	} else {
		requests, err = h.Service.List(r.Context(), programId, endpointId, jobId, rawSQL, domain, urlContains, urlMatch, includeSubdomains, orderBy1, asc1, orderBy2, asc2, orderBy3, asc3, orderBy4, asc4)
	}

	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*RequestList, len(requests))
	for i, req := range requests {
		response[i] = ToRequestList(req)
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
