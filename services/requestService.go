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

// extractContextSnippet extracts a contextual snippet around the matched text
func extractContextSnippet(text, searchQuery string, contextType string) string {
	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(searchQuery)

	// Find the position of the search query
	pos := strings.Index(lowerText, lowerQuery)
	if pos == -1 {
		return ""
	}

	// Split text into words
	words := strings.Fields(text)

	// Find the word position of the match
	wordPos := 0
	currentPos := 0
	for i, word := range words {
		if currentPos >= pos {
			wordPos = i
			break
		}
		currentPos += len(word) + 1 // +1 for space
	}

	// Calculate start and end positions for ~20 words with match at 7th/8th position
	startPos := wordPos - 7
	if startPos < 0 {
		startPos = 0
	}

	endPos := startPos + 20
	if endPos > len(words) {
		endPos = len(words)
		// Adjust start position if we're near the end
		if endPos-startPos < 20 && startPos > 0 {
			startPos = endPos - 20
			if startPos < 0 {
				startPos = 0
			}
		}
	}

	// Extract the snippet
	snippet := strings.Join(words[startPos:endPos], " ")

	// Add ellipsis if we're not at the beginning or end
	if startPos > 0 {
		snippet = "..." + snippet
	}
	if endPos < len(words) {
		snippet = snippet + "..."
	}

	return contextType + ": " + snippet
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
func (s *RequestService) SearchRequests(ctx context.Context, searchQuery string) ([]*models.MyRequest, map[int][]string, error) {
	var requests []*models.MyRequest
	searchResults := make(map[int][]string)

	// Search in request body, response body, headers, and URL
	query := s.DB.WithContext(ctx).Preload("Program").Preload("Endpoint").Preload("Notes").Preload("Attachments")

	searchPattern := "%" + searchQuery + "%"
	query = query.Where("url LIKE ? OR req_body LIKE ? OR res_body LIKE ? OR req_headers LIKE ? OR res_headers LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)

	if err := query.Find(&requests).Error; err != nil {
		return nil, nil, err
	}

	// Extract search results (matched strings) for each request
	for _, req := range requests {
		var reqSearchResults []string

		// Search in URL
		if strings.Contains(strings.ToLower(req.URL), strings.ToLower(searchQuery)) {
			snippet := extractContextSnippet(req.URL, searchQuery, "URL")
			if snippet != "" {
				reqSearchResults = append(reqSearchResults, snippet)
			}
		}

		// Search in request body
		if strings.Contains(strings.ToLower(req.ReqBody), strings.ToLower(searchQuery)) {
			snippet := extractContextSnippet(req.ReqBody, searchQuery, "Request Body")
			if snippet != "" {
				reqSearchResults = append(reqSearchResults, snippet)
			}
		}

		// Search in response body
		if strings.Contains(strings.ToLower(req.ResBody), strings.ToLower(searchQuery)) {
			snippet := extractContextSnippet(req.ResBody, searchQuery, "Response Body")
			if snippet != "" {
				reqSearchResults = append(reqSearchResults, snippet)
			}
		}

		// Search in headers
		if strings.Contains(strings.ToLower(req.ReqHeaders), strings.ToLower(searchQuery)) {
			snippet := extractContextSnippet(req.ReqHeaders, searchQuery, "Request Headers")
			if snippet != "" {
				reqSearchResults = append(reqSearchResults, snippet)
			}
		}

		if strings.Contains(strings.ToLower(req.ResHeaders), strings.ToLower(searchQuery)) {
			snippet := extractContextSnippet(req.ResHeaders, searchQuery, "Response Headers")
			if snippet != "" {
				reqSearchResults = append(reqSearchResults, snippet)
			}
		}

		searchResults[req.Id] = reqSearchResults
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
