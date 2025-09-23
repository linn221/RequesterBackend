package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"

	"github.com/linn221/RequesterBackend/har"
	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type ImportHarService struct {
	DB *gorm.DB
}

// ImportHAR processes a HAR file and creates an import job with associated requests and endpoints
func (s *ImportHarService) ImportHAR(ctx context.Context, file io.Reader, filename string, programId int, ignoredHeaders string) (int, error) {
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	// Create import job
	job := &models.ImportJob{
		ProgramId:      &programId,
		JobType:        "import_har",
		Title:          fmt.Sprintf("Import HAR: %s", filepath.Base(filename)),
		Progress:       0,
		Description:    fmt.Sprintf("Importing HAR file: %s", filename),
		IgnoredHeaders: ignoredHeaders,
	}

	if err := s.DB.WithContext(ctx).Create(job).Error; err != nil {
		return 0, fmt.Errorf("failed to create import job: %v", err)
	}

	// Update progress
	job.Progress = 10
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	// Parse HAR file
	requests, err := har.ParseHAR(fileContent, s.resHashFunc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse HAR file: %v", err)
	}

	// Update progress
	job.Progress = 30
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	// Group requests by domain and method to create endpoints
	endpointMap := make(map[string]*models.Endpoint)
	var endpoints []*models.Endpoint

	for _, req := range requests {
		// Validate that domain is set - this is critical
		if req.Domain == "" {
			// Try to extract domain from URL as fallback
			u, err := url.Parse(req.URL)
			if err == nil && u.Hostname() != "" {
				req.Domain = u.Hostname()
				log.Printf("HAR Import: Extracted domain '%s' from URL '%s' as fallback", req.Domain, req.URL)
			} else {
				// Skip requests without valid domain
				log.Printf("HAR Import: Skipping request with invalid URL '%s' - cannot extract domain", req.URL)
				continue
			}
		} else {
			log.Printf("HAR Import: Using domain '%s' for URL '%s'", req.Domain, req.URL)
		}

		// Create endpoint key (domain + method + path)
		u, err := url.Parse(req.URL)
		if err != nil {
			continue
		}

		path := u.Path
		if path == "" {
			path = "/"
		}

		endpointKey := fmt.Sprintf("%s:%s:%s", req.Domain, req.Method, path)

		if _, exists := endpointMap[endpointKey]; !exists {
			endpoint := &models.Endpoint{
				ProgramId:    programId,
				Method:       req.Method,
				Domain:       req.Domain,
				URI:          path,
				EndpointType: models.EndpointTypeAPI, // Default to API
				Note:         fmt.Sprintf("Auto-generated from HAR import: %s", filename),
			}

			endpointMap[endpointKey] = endpoint
			endpoints = append(endpoints, endpoint)
		}
	}

	// Save endpoints
	if len(endpoints) > 0 {
		if err := s.DB.WithContext(ctx).CreateInBatches(endpoints, 100).Error; err != nil {
			return 0, fmt.Errorf("failed to create endpoints: %v", err)
		}
	}

	// Update progress
	job.Progress = 60
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	// Create a map of endpoint key to endpoint ID for quick lookup
	endpointKeyToID := make(map[string]int)
	for _, endpoint := range endpoints {
		endpointKey := fmt.Sprintf("%s:%s:%s", endpoint.Domain, endpoint.Method, endpoint.URI)
		endpointKeyToID[endpointKey] = endpoint.Id
	}

	// Update requests with endpoint IDs, import job ID, and program ID
	for i := range requests {
		u, err := url.Parse(requests[i].URL)
		if err != nil {
			continue
		}

		path := u.Path
		if path == "" {
			path = "/"
		}

		endpointKey := fmt.Sprintf("%s:%s:%s", requests[i].Domain, requests[i].Method, path)
		if endpointID, exists := endpointKeyToID[endpointKey]; exists {
			requests[i].EndpointId = endpointID
		}
		requests[i].ImportJobId = job.Id
		requests[i].ProgramId = &programId // Set the program_id from the import form
	}

	// Save requests in batches
	if len(requests) > 0 {
		if err := s.DB.WithContext(ctx).CreateInBatches(requests, 100).Error; err != nil {
			return 0, fmt.Errorf("failed to create requests: %v", err)
		}
	}

	// Update progress to complete
	job.Progress = 100
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	return job.Id, nil
}

// resHashFunc is used by the HAR parser to generate request and response hashes
func (s *ImportHarService) resHashFunc(req *models.MyRequest) (string, string) {
	// Generate request text
	reqHeadersFromJSON, _ := models.HeaderSliceFromJSON(req.ReqHeaders)
	requestText := req.Method + " " + req.URL + " " + req.ReqBody + " " + reqHeadersFromJSON.EchoAll()

	// Generate response text
	resHeadersFromJSON, _ := models.HeaderSliceFromJSON(req.ResHeaders)
	responseText := fmt.Sprintf("%d %s", req.ResStatus, resHeadersFromJSON.EchoAll()) + req.ResBody

	return requestText, responseText
}

// GetImportJob retrieves an import job by ID
func (s *ImportHarService) GetImportJob(ctx context.Context, id int) (*models.ImportJob, error) {
	return first[models.ImportJob](s.DB.WithContext(ctx), id)
}

// ListImportJobs retrieves all import jobs
func (s *ImportHarService) ListImportJobs(ctx context.Context) ([]*models.ImportJob, error) {
	var jobs []*models.ImportJob
	if err := s.DB.WithContext(ctx).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
