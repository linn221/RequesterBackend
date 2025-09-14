package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/linn221/RequesterBackend/handlers"
	"github.com/linn221/RequesterBackend/services"
	"github.com/linn221/RequesterBackend/utils"
)

// RegisterRoutes registers all routes for the application
func (app *App) RegisterRoutes() *http.ServeMux {
	handler := handlers.Handler{}
	uploadDir := utils.GetEnv("UPLOAD_DIR", "./uploads")
	maxFileSize := utils.GetEnv("MAX_FILE_SIZE", "10485760") // 10MB in bytes
	allowedTypes := []string{
		"image/",
		"application/pdf",
		"text/",
		"application/json",
		"application/octet-stream",
	}

	mux := http.NewServeMux()

	// Notes
	noteService := services.NoteService{
		DB: app.DB,
	}
	noteHandler := handlers.NoteHandler{
		Service: &noteService,
	}
	mux.HandleFunc("POST /notes", noteHandler.Create)
	mux.HandleFunc("GET /notes", noteHandler.List)
	mux.HandleFunc("GET /notes/{id}", noteHandler.Get)
	mux.HandleFunc("DELETE /notes/{id}", noteHandler.Delete)
	mux.HandleFunc("PATCH /notes/{id}", noteHandler.Update)

	// Attachments
	attachmentService := services.AttachmentService{
		DB:              app.DB,
		UploadDirectory: uploadDir,
		MaxFileSize:     parseMaxFileSize(maxFileSize),
		AllowedTypes:    allowedTypes,
		ValidTypes:      []string{"programs", "endpoints", "requests"},
	}
	attachmentHandler := handlers.AttachmentHandler{
		AttachmentService: &attachmentService,
	}
	mux.HandleFunc("POST /attachments", attachmentHandler.UploadAttachment)
	mux.HandleFunc("DELETE /attachments", attachmentHandler.DeleteAttachment)
	mux.HandleFunc("GET /attachments/{id}", attachmentHandler.GetAttachment)

	// Programs
	programService := services.ProgramService{
		DB: app.DB,
	}
	programHandler := handlers.ProgramHandler{
		Service: &programService,
	}
	mux.HandleFunc("POST /programs", programHandler.Create)
	mux.HandleFunc("GET /programs", programHandler.List)
	mux.HandleFunc("GET /programs/{id}", programHandler.Get)
	mux.HandleFunc("PUT /programs/{id}", programHandler.Update)
	mux.HandleFunc("DELETE /programs/{id}", programHandler.Delete)

	// Endpoints
	endpointService := services.EndpointService{
		DB: app.DB,
	}
	endpointHandler := handlers.EndpointHandler{
		Service: &endpointService,
	}
	mux.HandleFunc("POST /endpoints", endpointHandler.Create)
	mux.HandleFunc("GET /endpoints", endpointHandler.List)
	mux.HandleFunc("GET /endpoints/{id}", endpointHandler.Get)
	mux.HandleFunc("PUT /endpoints/{id}", endpointHandler.Update)
	mux.HandleFunc("DELETE /endpoints/{id}", endpointHandler.Delete)

	// Requests
	mux.HandleFunc("GET /requests", handler.ListRequests)
	mux.HandleFunc("GET /requests/{id}", handler.GetRequest)

	// Import HAR
	importHarService := services.ImportHarService{
		DB: app.DB,
	}
	importHarHandler := handlers.ImportHarHandler{
		Service: &importHarService,
	}
	mux.HandleFunc("POST /import_har", importHarHandler.ImportHAR)

	// Jobs
	jobHandler := handlers.JobHandler{
		Service: &importHarService,
	}
	mux.HandleFunc("GET /jobs", jobHandler.ListJobs)
	return mux
}

// parseMaxFileSize parses the max file size from string to int64
func parseMaxFileSize(sizeStr string) int64 {
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		log.Printf("Invalid MAX_FILE_SIZE value '%s', using default 10MB", sizeStr)
		return 10 * 1024 * 1024 // 10MB default
	}
	return size
}
