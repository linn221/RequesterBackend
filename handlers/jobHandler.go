package handlers

import (
	"net/http"

	"github.com/linn221/RequesterBackend/models"
	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

type JobHandler struct {
	Service *services.ImportHarService
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.Service.ListImportJobs(r.Context())
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	response := make([]*Job, len(jobs))
	for i, job := range jobs {
		response[i] = ToJob(job)
	}

	utils.OkJson(w, response)
}

func ToJob(job *models.ImportJob) *Job {
	return &Job{
		Id:          job.Id,
		JobType:     job.JobType,
		Title:       job.Title,
		Progress:    job.Progress,
		CreatedAt:   job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Description: job.Description,
	}
}
