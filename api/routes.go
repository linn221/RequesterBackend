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
		ValidTypes:      []string{"programs", "endpoints", "requests", "vulns"},
	}
	attachmentHandler := handlers.AttachmentHandler{
		AttachmentService: &attachmentService,
	}
	mux.HandleFunc("POST /attachments", attachmentHandler.UploadAttachment)
	mux.HandleFunc("DELETE /attachments", attachmentHandler.DeleteAttachment)
	mux.HandleFunc("GET /attachments/{id}", attachmentHandler.GetAttachment)

	// Images
	imageService := services.ImageService{
		DB:              app.DB,
		UploadDirectory: uploadDir,
		ValidTypes:      []string{"programs", "endpoints", "requests", "vulns"},
	}
	imageHandler := handlers.ImageHandler{
		ImageService: &imageService,
	}
	mux.HandleFunc("POST /images", imageHandler.UploadImages)
	mux.HandleFunc("DELETE /images", imageHandler.DeleteImage)
	mux.HandleFunc("GET /images/{id}", imageHandler.GetImage)
	mux.HandleFunc("GET /images/file/{filename}", imageHandler.ServeImage)

	// Vulnerabilities
	vulnService := services.VulnService{
		DB: app.DB,
	}
	vulnHandler := handlers.VulnHandler{
		Service: &vulnService,
	}
	mux.HandleFunc("POST /vulns", vulnHandler.Create)
	mux.HandleFunc("GET /vulns", vulnHandler.List)
	mux.HandleFunc("GET /vulns/{id}", vulnHandler.Get)
	mux.HandleFunc("GET /vulns/slug/{slug}", vulnHandler.GetBySlug)
	mux.HandleFunc("PUT /vulns/{id}", vulnHandler.Update)
	mux.HandleFunc("DELETE /vulns/{id}", vulnHandler.Delete)

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
	requestService := services.RequestService{
		DB: app.DB,
	}
	requestHandler := handlers.RequestHandler{
		Service: &requestService,
	}
	mux.HandleFunc("GET /requests", requestHandler.List)
	mux.HandleFunc("GET /requests/{id}", requestHandler.Get)

	// Import HAR
	importHarService := services.ImportHarService{
		DB: app.DB,
	}
	importHarHandler := handlers.ImportHarHandler{
		Service: &importHarService,
	}
	mux.HandleFunc("POST /import_har", importHarHandler.ImportHAR)

	// Import Burp XML
	importBurpService := services.ImportBurpService{
		DB: app.DB,
	}
	importBurpHandler := handlers.ImportBurpHandler{
		Service: &importBurpService,
	}
	mux.HandleFunc("POST /import_burp_xml", importBurpHandler.ImportBurpXML)

	// Jobs
	jobHandler := handlers.JobHandler{
		Service: &importHarService,
	}
	mux.HandleFunc("GET /jobs", jobHandler.ListJobs)

	// Add Swagger UI routes for development
	// To disable in production, comment out the following line:
	addSwaggerRoutes(mux)

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

// addSwaggerRoutes adds Swagger UI routes for development
// To disable in production, comment out the call to this function in RegisterRoutes()
func addSwaggerRoutes(mux *http.ServeMux) {
	// Serve OpenAPI spec
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "./openapi.yaml")
	})

	// Create a file server for Swagger UI static files
	swaggerFS := http.FileServer(http.Dir("./swagger/"))

	// Handle /swagger route - redirect to /swagger/
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// Handle /swagger/ and all sub-paths
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", swaggerFS))
}
