package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type RequestService struct {
	DB *gorm.DB
}

// List retrieves requests with filtering and search
func (s *RequestService) List(ctx context.Context, programId, endpointId, jobId *int, rawSQL, orderBy string, asc bool) ([]*models.MyRequest, error) {
	var requests []*models.MyRequest
	query := s.DB.WithContext(ctx).Preload("Program").Preload("Endpoint").Preload("Notes").Preload("Attachments")

	// Apply filters
	if programId != nil {
		query = query.Where("program_id = ?", *programId)
	}
	if endpointId != nil {
		query = query.Where("endpoint_id = ?", *endpointId)
	}
	if jobId != nil {
		query = query.Where("import_job_id = ?", *jobId)
	}

	// Apply raw SQL filter if provided
	if rawSQL != "" {
		query = query.Where(rawSQL)
	}

	// Apply ordering
	if orderBy != "" {
		orderDirection := "ASC"
		if !asc {
			orderDirection = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDirection))
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// Get retrieves a request by Id
func (s *RequestService) Get(ctx context.Context, id int) (*models.MyRequest, error) {
	var request models.MyRequest
	if err := s.DB.WithContext(ctx).Preload("Program").Preload("Endpoint").Preload("Notes").Preload("Attachments").First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// SearchRequests searches for requests based on query string
func (s *RequestService) SearchRequests(ctx context.Context, searchQuery string) ([]*models.MyRequest, []string, error) {
	var requests []*models.MyRequest
	var searchResults []string

	// Search in request body, response body, headers, and URL
	query := s.DB.WithContext(ctx).Preload("Program").Preload("Endpoint").Preload("Notes").Preload("Attachments")

	searchPattern := "%" + searchQuery + "%"
	query = query.Where("url LIKE ? OR req_body LIKE ? OR res_body LIKE ? OR req_headers LIKE ? OR res_headers LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)

	if err := query.Find(&requests).Error; err != nil {
		return nil, nil, err
	}

	// Extract search results (matched strings)
	for _, req := range requests {
		// Search in URL
		if strings.Contains(strings.ToLower(req.URL), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, "URL: "+req.URL)
		}

		// Search in request body
		if strings.Contains(strings.ToLower(req.ReqBody), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, "Request Body: "+req.ReqBody)
		}

		// Search in response body
		if strings.Contains(strings.ToLower(req.ResBody), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, "Response Body: "+req.ResBody)
		}

		// Search in headers
		if strings.Contains(strings.ToLower(req.ReqHeaders), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, "Request Headers: "+req.ReqHeaders)
		}

		if strings.Contains(strings.ToLower(req.ResHeaders), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, "Response Headers: "+req.ResHeaders)
		}
	}

	return requests, searchResults, nil
}

// ParseRequestHeaders parses JSON headers string to map
func ParseRequestHeaders(headersJSON string) (map[string]interface{}, error) {
	var headers map[string]interface{}
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return nil, err
	}
	return headers, nil
}

// ParseResponseHeaders parses JSON headers string to map
func ParseResponseHeaders(headersJSON string) (map[string]interface{}, error) {
	var headers map[string]interface{}
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return nil, err
	}
	return headers, nil
}

// ParseRequestBody parses JSON body string to interface
func ParseRequestBody(bodyJSON string) (interface{}, error) {
	var body interface{}
	if err := json.Unmarshal([]byte(bodyJSON), &body); err != nil {
		// If it's not JSON, return as string
		return bodyJSON, nil
	}
	return body, nil
}

// ParseResponseBody parses JSON body string to interface
func ParseResponseBody(bodyJSON string) (interface{}, error) {
	var body interface{}
	if err := json.Unmarshal([]byte(bodyJSON), &body); err != nil {
		// If it's not JSON, return as string
		return bodyJSON, nil
	}
	return body, nil
}
