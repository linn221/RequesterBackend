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

// buildOrderClauses builds ORDER BY clauses from multiple ordering parameters
func (s *RequestService) buildOrderClauses(orderBy1 string, asc1 bool, orderBy2 string, asc2 bool, orderBy3 string, asc3 bool, orderBy4 string, asc4 bool) string {
	var clauses []string

	// Helper function to map order_by parameter to database column
	mapToColumn := func(orderBy string) string {
		switch orderBy {
		case "method":
			return "method"
		case "content_type":
			// content_type is extracted from response headers, we'll need to handle this differently
			// For now, skip ordering by content_type as it requires a more complex query
			return ""
		case "size":
			return "resp_size"
		case "latency":
			return "latency_ms"
		case "url":
			return "url"
		case "sequence_number":
			return "sequence"
		default:
			// Invalid order_by parameter, skip ordering
			return ""
		}
	}

	// Helper function to build a single order clause
	buildClause := func(orderBy string, asc bool) string {
		column := mapToColumn(orderBy)
		if column == "" {
			return ""
		}
		direction := "ASC"
		if !asc {
			direction = "DESC"
		}
		return fmt.Sprintf("%s %s", column, direction)
	}

	// Build clauses for each level (orderBy1 has highest priority)
	if orderBy1 != "" {
		if clause := buildClause(orderBy1, asc1); clause != "" {
			clauses = append(clauses, clause)
		}
	}
	if orderBy2 != "" {
		if clause := buildClause(orderBy2, asc2); clause != "" {
			clauses = append(clauses, clause)
		}
	}
	if orderBy3 != "" {
		if clause := buildClause(orderBy3, asc3); clause != "" {
			clauses = append(clauses, clause)
		}
	}
	if orderBy4 != "" {
		if clause := buildClause(orderBy4, asc4); clause != "" {
			clauses = append(clauses, clause)
		}
	}

	// Join all clauses with commas
	return strings.Join(clauses, ", ")
}

// List retrieves requests with filtering and search
func (s *RequestService) List(ctx context.Context, programId, endpointId, jobId *int, rawSQL, orderBy1 string, asc1 bool, orderBy2 string, asc2 bool, orderBy3 string, asc3 bool, orderBy4 string, asc4 bool) ([]*models.MyRequest, error) {
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

	// Apply multi-level ordering
	orderClauses := s.buildOrderClauses(orderBy1, asc1, orderBy2, asc2, orderBy3, asc3, orderBy4, asc4)
	if len(orderClauses) > 0 {
		query = query.Order(orderClauses)
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
func (s *RequestService) SearchRequests(ctx context.Context, searchQuery, orderBy1 string, asc1 bool, orderBy2 string, asc2 bool, orderBy3 string, asc3 bool, orderBy4 string, asc4 bool) ([]*models.MyRequest, error) {
	var requests []*models.MyRequest

	// Search in request body, response body, headers, and URL
	query := s.DB.WithContext(ctx).Preload("Program").Preload("Endpoint").Preload("Notes").Preload("Attachments")

	searchPattern := "%" + searchQuery + "%"
	query = query.Where("url LIKE ? OR req_body LIKE ? OR res_body LIKE ? OR req_headers LIKE ? OR res_headers LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)

	// Apply multi-level ordering
	orderClauses := s.buildOrderClauses(orderBy1, asc1, orderBy2, asc2, orderBy3, asc3, orderBy4, asc4)
	if len(orderClauses) > 0 {
		query = query.Order(orderClauses)
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
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
