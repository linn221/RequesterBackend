package services

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/linn221/RequesterBackend/models"
	"github.com/linn221/RequesterBackend/utils"
	"gorm.io/gorm"
)

type ImportBurpService struct {
	DB *gorm.DB
}

// filterHeaders filters out ignored headers from a HeaderSlice
func (s *ImportBurpService) filterHeaders(headers []models.Header, ignoredHeaders string) []models.Header {
	if ignoredHeaders == "" {
		return headers
	}

	ignoredList := strings.Split(ignoredHeaders, ",")
	ignoredMap := make(map[string]struct{})
	for _, h := range ignoredList {
		ignoredMap[strings.ToLower(strings.TrimSpace(h))] = struct{}{}
	}

	var filtered []models.Header
	for _, h := range headers {
		if _, ignored := ignoredMap[strings.ToLower(h.Name)]; !ignored {
			filtered = append(filtered, h)
		}
	}

	return filtered
}

// resHashFunc is used to generate request and response hashes
func (s *ImportBurpService) resHashFunc(req *models.MyRequest) (string, string) {
	// Generate request text
	reqHeadersFromJSON, _ := models.HeaderSliceFromJSON(req.ReqHeaders)
	requestText := req.Method + " " + req.URL + " " + req.ReqBody + " " + reqHeadersFromJSON.EchoAll()

	// Generate response text
	resHeadersFromJSON, _ := models.HeaderSliceFromJSON(req.ResHeaders)
	responseText := fmt.Sprintf("%d %s", req.ResStatus, resHeadersFromJSON.EchoAll()) + req.ResBody

	return requestText, responseText
}

// BurpXML represents the structure of a Burp Suite XML file
type BurpXML struct {
	Items []BurpItem `xml:"item"`
}

type BurpItem struct {
	Time      string `xml:"time"`
	URL       string `xml:"url"`
	Host      string `xml:"host"`
	Port      int    `xml:"port"`
	Protocol  string `xml:"protocol"`
	Method    string `xml:"method"`
	Path      string `xml:"path"`
	Extension string `xml:"extension"`
	Request   struct {
		Base64 string `xml:",chardata"`
	} `xml:"request"`
	Response struct {
		Base64 string `xml:",chardata"`
	} `xml:"response"`
	Comment string `xml:"comment"`
}

// ImportBurpXML processes a Burp XML file and creates an import job with associated requests and endpoints
func (s *ImportBurpService) ImportBurpXML(ctx context.Context, file io.Reader, filename string, programId int, ignoredHeaders string) (int, error) {
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	// Create import job
	job := &models.ImportJob{
		ProgramId:      &programId,
		JobType:        "import_burp_xml",
		Title:          fmt.Sprintf("Import Burp XML: %s", filepath.Base(filename)),
		Progress:       0,
		Description:    fmt.Sprintf("Importing Burp XML file: %s", filename),
		IgnoredHeaders: ignoredHeaders,
	}

	if err := s.DB.WithContext(ctx).Create(job).Error; err != nil {
		return 0, fmt.Errorf("failed to create import job: %v", err)
	}

	// Update progress
	job.Progress = 10
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	// Parse Burp XML file
	var burpXML BurpXML
	if err := xml.Unmarshal(fileContent, &burpXML); err != nil {
		return 0, fmt.Errorf("failed to parse Burp XML file: %v", err)
	}

	// Update progress
	job.Progress = 30
	s.DB.WithContext(ctx).Model(job).Update("Progress", job.Progress)

	// Group requests by domain and method to create endpoints
	endpointMap := make(map[string]*models.Endpoint)
	var endpoints []*models.Endpoint

	// Process each item in the Burp XML
	for _, item := range burpXML.Items {
		// Decode base64 request
		requestData, err := base64.StdEncoding.DecodeString(item.Request.Base64)
		if err != nil {
			continue // Skip items with invalid base64 data
		}

		// Parse the request to extract headers and body
		requestText := string(requestData)
		requestLines := strings.Split(requestText, "\n")
		if len(requestLines) < 2 {
			continue // Skip malformed requests
		}

		// Extract method and URI from first line
		firstLine := requestLines[0]
		parts := strings.Split(firstLine, " ")
		if len(parts) < 3 {
			continue
		}

		method := parts[0]
		uri := parts[1]

		// Construct full URL using host, port, and protocol from Burp XML
		var scheme string
		if item.Protocol == "https" {
			scheme = "https"
		} else {
			scheme = "http"
		}

		var port string
		if item.Port != 0 && item.Port != 80 && item.Port != 443 {
			port = fmt.Sprintf(":%d", item.Port)
		}

		// Use the URL from the XML if available, otherwise construct it
		var requestURL string
		if item.URL != "" {
			requestURL = item.URL
		} else {
			requestURL = fmt.Sprintf("%s://%s%s%s", scheme, item.Host, port, uri)
		}

		// Parse headers
		var headers []models.Header
		bodyStart := 0
		for i, line := range requestLines[1:] {
			if line == "" {
				bodyStart = i + 2
				break
			}
			if strings.Contains(line, ":") {
				headerParts := strings.SplitN(line, ":", 2)
				if len(headerParts) == 2 {
					headers = append(headers, models.Header{
						Name:  strings.TrimSpace(headerParts[0]),
						Value: strings.TrimSpace(headerParts[1]),
					})
				}
			}
		}

		// Extract body (not used in endpoint creation, but needed for parsing)
		if bodyStart < len(requestLines) {
			_ = strings.Join(requestLines[bodyStart:], "\n")
		}

		// Parse URL to get domain and path
		parsedURL, err := url.Parse(requestURL)
		if err != nil {
			continue
		}

		domain := parsedURL.Hostname() // Use Hostname() to get just the domain without port
		path := parsedURL.Path
		if path == "" {
			path = "/"
		}

		// Create endpoint key (domain + method + path)
		endpointKey := fmt.Sprintf("%s:%s:%s", domain, method, path)

		if _, exists := endpointMap[endpointKey]; !exists {
			endpoint := &models.Endpoint{
				ProgramId:    programId,
				Method:       method,
				Domain:       domain,
				URI:          path,
				EndpointType: models.EndpointTypeAPI, // Default to API
				Note:         fmt.Sprintf("Auto-generated from Burp XML import: %s", filename),
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

	// Process requests and create MyRequest objects
	var requests []*models.MyRequest

	for i, item := range burpXML.Items {
		// Decode base64 request and response
		requestData, err := base64.StdEncoding.DecodeString(item.Request.Base64)
		if err != nil {
			continue
		}

		responseData, err := base64.StdEncoding.DecodeString(item.Response.Base64)
		if err != nil {
			continue
		}

		// Parse the request
		requestText := string(requestData)
		requestLines := strings.Split(requestText, "\n")
		if len(requestLines) < 2 {
			continue
		}

		// Extract method and URI from first line
		firstLine := requestLines[0]
		parts := strings.Split(firstLine, " ")
		if len(parts) < 3 {
			continue
		}

		method := parts[0]
		uri := parts[1]

		// Construct full URL using host, port, and protocol from Burp XML
		var scheme string
		if item.Protocol == "https" {
			scheme = "https"
		} else {
			scheme = "http"
		}

		var port string
		if item.Port != 0 && item.Port != 80 && item.Port != 443 {
			port = fmt.Sprintf(":%d", item.Port)
		}

		// Use the URL from the XML if available, otherwise construct it
		var requestURL string
		if item.URL != "" {
			requestURL = item.URL
		} else {
			requestURL = fmt.Sprintf("%s://%s%s%s", scheme, item.Host, port, uri)
		}

		// Parse headers
		var headers []models.Header
		bodyStart := 0
		for i, line := range requestLines[1:] {
			if line == "" {
				bodyStart = i + 2
				break
			}
			if strings.Contains(line, ":") {
				headerParts := strings.SplitN(line, ":", 2)
				if len(headerParts) == 2 {
					headers = append(headers, models.Header{
						Name:  strings.TrimSpace(headerParts[0]),
						Value: strings.TrimSpace(headerParts[1]),
					})
				}
			}
		}

		// Extract body
		var body string
		if bodyStart < len(requestLines) {
			body = strings.Join(requestLines[bodyStart:], "\n")
		}

		// Parse URL to get domain and path
		parsedURL, err := url.Parse(requestURL)
		if err != nil {
			continue
		}

		domain := parsedURL.Hostname() // Use Hostname() instead of Host to get just the domain
		path := parsedURL.Path
		if path == "" {
			path = "/"
		}

		// Parse response to extract status code, headers, and body
		responseText := string(responseData)
		responseLines := strings.Split(responseText, "\n")
		var statusCode int = 200 // Default status code
		var responseHeaders []models.Header
		var responseBody string

		if len(responseLines) > 0 {
			// Extract status code from first line
			statusLine := responseLines[0]
			if strings.HasPrefix(statusLine, "HTTP/") {
				statusParts := strings.Split(statusLine, " ")
				if len(statusParts) >= 2 {
					if code, err := strconv.Atoi(statusParts[1]); err == nil {
						statusCode = code
					}
				}
			}

			// Parse response headers - look for the empty line that separates headers from body
			responseBodyStart := 0
			for i, line := range responseLines[1:] {
				// Check for empty line that separates headers from body
				if strings.TrimSpace(line) == "" {
					responseBodyStart = i + 2 // +2 because we're iterating from responseLines[1:]
					break
				}
				// Parse header lines (format: "Header-Name: Header-Value")
				if strings.Contains(line, ":") {
					headerParts := strings.SplitN(line, ":", 2)
					if len(headerParts) == 2 {
						responseHeaders = append(responseHeaders, models.Header{
							Name:  strings.TrimSpace(headerParts[0]),
							Value: strings.TrimSpace(headerParts[1]),
						})
					}
				}
			}

			// Extract response body - everything after the empty line
			if responseBodyStart > 0 && responseBodyStart < len(responseLines) {
				responseBody = strings.Join(responseLines[responseBodyStart:], "\n")
			}
		}

		// Filter headers based on ignored headers
		filteredReqHeaders := s.filterHeaders(headers, ignoredHeaders)
		filteredResHeaders := s.filterHeaders(responseHeaders, ignoredHeaders)

		// Parse time for latency calculation
		var latencyMs int64 = 0
		var requestTime string
		if item.Time != "" {
			// Try to parse the time - Burp uses various formats
			if t, err := time.Parse("Mon Jan 2 15:04:05 GMT-07:00 2006", item.Time); err == nil {
				requestTime = t.Format(time.RFC3339)
				// Calculate latency (this is approximate since we don't have response time)
				// We'll use 0 for now, but this could be enhanced if needed
				latencyMs = 0
			} else {
				// Fallback to current time if parsing fails
				requestTime = time.Now().Format(time.RFC3339)
			}
		} else {
			requestTime = time.Now().Format(time.RFC3339)
		}

		// Create endpoint key
		endpointKey := fmt.Sprintf("%s:%s:%s", domain, method, path)

		// Create MyRequest object with all fields
		request := &models.MyRequest{
			Sequence:    i + 1, // Add sequence number
			URL:         requestURL,
			Method:      method,
			Domain:      domain,
			ReqBody:     body,
			ResBody:     responseBody,
			ResStatus:   statusCode,
			RespSize:    len(responseBody), // Add response size
			LatencyMs:   latencyMs,         // Add latency
			RequestTime: requestTime,       // Add request time
			ProgramId:   &programId,
			ImportJobId: job.Id,
		}

		// Set endpoint ID if found
		if endpointID, exists := endpointKeyToID[endpointKey]; exists {
			request.EndpointId = endpointID
		}

		// Convert filtered headers to JSON and truncate if too long
		reqHeadersSlice := models.HeaderSlice(filteredReqHeaders)
		resHeadersSlice := models.HeaderSlice(filteredResHeaders)
		reqHeadersJSON, _ := reqHeadersSlice.ToJSON()
		resHeadersJSON, _ := resHeadersSlice.ToJSON()

		// Truncate headers if they exceed MySQL TEXT column limit (65,535 characters)
		const maxTextLength = 65000 // Leave some buffer
		if len(reqHeadersJSON) > maxTextLength {
			reqHeadersJSON = reqHeadersJSON[:maxTextLength]
		}
		if len(resHeadersJSON) > maxTextLength {
			resHeadersJSON = resHeadersJSON[:maxTextLength]
		}

		request.ReqHeaders = reqHeadersJSON
		request.ResHeaders = resHeadersJSON

		// Calculate hashes using the resHashFunc
		reqText, resText := s.resHashFunc(request)
		request.ReqHash = utils.HashString(reqText)
		request.ResHash = utils.HashString(resText)
		request.ResBodyHash = utils.HashString(request.ResBody)
		reqHeadersFromJSON, _ := models.HeaderSliceFromJSON(request.ReqHeaders)
		request.ReqHash1 = utils.HashString(request.Method + " " + request.URL + " " + request.ReqBody + " " + reqHeadersFromJSON.EchoAll())

		requests = append(requests, request)
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

// GetImportJob retrieves an import job by ID
func (s *ImportBurpService) GetImportJob(ctx context.Context, id int) (*models.ImportJob, error) {
	return first[models.ImportJob](s.DB.WithContext(ctx), id)
}

// ListImportJobs retrieves all import jobs
func (s *ImportBurpService) ListImportJobs(ctx context.Context) ([]*models.ImportJob, error) {
	var jobs []*models.ImportJob
	if err := s.DB.WithContext(ctx).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
