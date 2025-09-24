package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/linn221/RequesterBackend/models"
	"github.com/linn221/RequesterBackend/services"
)

// ===== Notes =====
type CreateNoteRequest struct {
	ReferenceType string `json:"reference_type" validate:"required,oneof=programs endpoints requests vulns"`
	ReferenceId   int    `json:"reference_id" validate:"required"`
	Value         string `json:"value" validate:"required"`
}

type Note struct {
	Id            int    `json:"id"`
	ReferenceType string `json:"reference_type"`
	ReferenceId   int    `json:"reference_id"`
	Value         string `json:"value"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type NoteListing struct {
	Id        int    `json:"id"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NoteDetail struct {
	Id            int      `json:"id"`
	ReferenceType string   `json:"reference_type"`
	ReferenceId   int      `json:"reference_id"`
	Value         string   `json:"value"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	Tags          []TagDTO `json:"tags"`
}

func ToNote(note *models.Note) *Note {
	return &Note{
		Id:            note.Id,
		ReferenceType: note.ReferenceType,
		ReferenceId:   note.ReferenceID,
		Value:         note.Value,
		CreatedAt:     note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToNoteListing(note *models.Note) *NoteListing {
	return &NoteListing{
		Id:        note.Id,
		Value:     note.Value,
		CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToNoteDetail(note *models.Note) *NoteDetail {
	tags := make([]TagDTO, len(note.Tags))
	for i, tag := range note.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &NoteDetail{
		Id:            note.Id,
		ReferenceType: note.ReferenceType,
		ReferenceId:   note.ReferenceID,
		Value:         note.Value,
		CreatedAt:     note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Tags:          tags,
	}
}

// ===== Attachments =====
type Attachment struct {
	Id           int    `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_filename"`
	URL          string `json:"url"`
}

type UploadAttachmentRequest struct {
	ReferenceType string `json:"reference_type" validate:"required,oneof=programs endpoints requests vulns"`
	ReferenceId   int    `json:"reference_id" validate:"required"`
	// file is handled separately as multipart
}

// ===== Images =====
type Image struct {
	Id           int    `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_filename"`
	URL          string `json:"url"`
}

type UploadImageRequest struct {
	ReferenceType string `json:"reference_type" validate:"required,oneof=programs endpoints requests vulns"`
	ReferenceId   int    `json:"reference_id" validate:"required"`
	// files are handled separately as multipart
}

// ===== Vulnerabilities =====
type VulnInput struct {
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
	ParentId *int   `json:"parent_id"`
}

func (input *VulnInput) ToModel() *models.Vuln {
	return &models.Vuln{
		Title:    input.Title,
		Body:     input.Body,
		ParentId: input.ParentId,
	}
}

type VulnList struct {
	Id         int      `json:"id"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	ParentId   *int     `json:"parent_id"`
	ParentName *string  `json:"parent_name"`
	Tags       []TagDTO `json:"tags"`
}

func ToVulnList(vuln *models.Vuln) *VulnList {
	var parentName *string
	if vuln.Parent != nil {
		parentName = &vuln.Parent.Title
	}

	tags := make([]TagDTO, len(vuln.Tags))
	for i, tag := range vuln.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &VulnList{
		Id:         vuln.Id,
		Title:      vuln.Title,
		Slug:       vuln.Slug,
		ParentId:   vuln.ParentId,
		ParentName: parentName,
		Tags:       tags,
	}
}

type VulnDetail struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	Slug        string        `json:"slug"`
	ParentId    *int          `json:"parent_id"`
	ParentVuln  *string       `json:"parent_vuln"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Notes       []NoteListing `json:"notes"`
	Attachments []Attachment  `json:"attachments"`
	Images      []Image       `json:"images"`
	Tags        []TagDTO      `json:"tags"`
}

func ToVulnDetail(vuln *models.Vuln) *VulnDetail {
	var parentVuln *string
	if vuln.Parent != nil {
		parentVuln = &vuln.Parent.Title
	}

	// Convert notes to NoteListing
	notes := make([]NoteListing, len(vuln.Notes))
	for i, note := range vuln.Notes {
		notes[i] = *ToNoteListing(&note)
	}

	// Convert attachments to Attachment
	attachments := make([]Attachment, len(vuln.Attachments))
	for i, attachment := range vuln.Attachments {
		attachments[i] = Attachment{
			Id:           attachment.Id,
			Filename:     attachment.Filename,
			OriginalName: attachment.OriginalName,
			URL:          attachment.GetURL(),
		}
	}

	// Convert images to Image
	images := make([]Image, len(vuln.Images))
	for i, image := range vuln.Images {
		images[i] = Image{
			Id:           image.Id,
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.GetURL(),
		}
	}

	// Convert tags to TagDTO
	tags := make([]TagDTO, len(vuln.Tags))
	for i, tag := range vuln.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &VulnDetail{
		Id:          vuln.Id,
		Title:       vuln.Title,
		Body:        vuln.Body,
		Slug:        vuln.Slug,
		ParentId:    vuln.ParentId,
		ParentVuln:  parentVuln,
		CreatedAt:   vuln.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   vuln.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Notes:       notes,
		Attachments: attachments,
		Images:      images,
		Tags:        tags,
	}
}

// ===== Programs =====
type ProgramInput struct {
	Name    string `json:"name" validate:"required"`
	URL     string `json:"url" validate:"required,url"`
	Scope   string `json:"scope"`
	Domains string `json:"domains"`
	Note    string `json:"note"`
}

func (input *ProgramInput) ToModel() *models.Program {
	return &models.Program{
		Name:    input.Name,
		URL:     input.URL,
		Scope:   input.Scope,
		Domains: input.Domains,
		Note:    input.Note,
	}
}

type ProgramList struct {
	Id   int      `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []TagDTO `json:"tags"`
}

func ToProgramList(program *models.Program) *ProgramList {
	tags := make([]TagDTO, len(program.Tags))
	for i, tag := range program.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &ProgramList{
		Id:   program.Id,
		Name: program.Name,
		URL:  program.URL,
		Tags: tags,
	}
}

type ProgramDetail struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Scope       string        `json:"scope"`
	Domains     string        `json:"domains"`
	Note        string        `json:"note"`
	Notes       []NoteListing `json:"notes"`
	Attachments []Attachment  `json:"attachments"`
	Images      []Image       `json:"images"`
	Tags        []TagDTO      `json:"tags"`
}

func ToProgramDetail(program *models.Program) *ProgramDetail {
	notes := make([]NoteListing, len(program.Notes))
	for i, note := range program.Notes {
		notePtr := &note
		notes[i] = *ToNoteListing(notePtr)
	}

	attachments := make([]Attachment, len(program.Attachments))
	for i, attachment := range program.Attachments {
		attachments[i] = Attachment{
			Id:           attachment.Id,
			Filename:     attachment.Filename,
			OriginalName: attachment.OriginalName,
			URL:          attachment.GetURL(),
		}
	}

	images := make([]Image, len(program.Images))
	for i, image := range program.Images {
		images[i] = Image{
			Id:           image.Id,
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.GetURL(),
		}
	}

	tags := make([]TagDTO, len(program.Tags))
	for i, tag := range program.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &ProgramDetail{
		Id:          program.Id,
		Name:        program.Name,
		URL:         program.URL,
		Scope:       program.Scope,
		Domains:     program.Domains,
		Note:        program.Note,
		Notes:       notes,
		Attachments: attachments,
		Images:      images,
		Tags:        tags,
	}
}

// ===== Endpoints =====
type EndpointInput struct {
	Domain       string `json:"domain" validate:"required"`
	ProgramId    int    `json:"program_id" validate:"required"`
	Method       string `json:"method" validate:"required"`
	URI          string `json:"uri" validate:"required"`
	EndpointType string `json:"endpoint_type" validate:"required,oneof=web api"`
	Description  string `json:"description" validate:"required"`
}

func (input *EndpointInput) ToModel() *models.Endpoint {
	return &models.Endpoint{
		ProgramId:    input.ProgramId,
		Method:       input.Method,
		Domain:       input.Domain,
		URI:          input.URI,
		EndpointType: models.EndpointType(input.EndpointType),
		Note:         input.Description,
	}
}

type EndpointList struct {
	Id           int      `json:"id"`
	ProgramId    int      `json:"program_id"`
	ProgramName  string   `json:"program_name"`
	Domain       string   `json:"domain"`
	URI          string   `json:"uri"`
	Method       string   `json:"method"`
	EndpointType string   `json:"endpoint_type"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	Text         string   `json:"text"`
	Tags         []TagDTO `json:"tags"`
}

func ToEndpointList(endpoint *models.Endpoint) *EndpointList {
	programName := ""
	if endpoint.Program != nil {
		programName = endpoint.Program.Name
	}

	// Concatenate all information into text field
	text := fmt.Sprintf("ID: %d\nProgram ID: %d\nProgram Name: %s\nDomain: %s\nURI: %s\nMethod: %s\nEndpoint Type: %s\nDescription: %s\nCreated At: %s\nUpdated At: %s",
		endpoint.Id,
		endpoint.ProgramId,
		programName,
		endpoint.Domain,
		endpoint.URI,
		endpoint.Method,
		string(endpoint.EndpointType),
		endpoint.Note,
		endpoint.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		endpoint.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	)

	// Add notes to text
	if len(endpoint.Notes) > 0 {
		text += "\nNotes:\n"
		for _, note := range endpoint.Notes {
			text += fmt.Sprintf("- %s\n", note.Value)
		}
	}

	// Add attachments to text
	if len(endpoint.Attachments) > 0 {
		text += "\nAttachments:\n"
		for _, attachment := range endpoint.Attachments {
			text += fmt.Sprintf("- %s (%s)\n", attachment.OriginalName, attachment.Filename)
		}
	}

	tags := make([]TagDTO, len(endpoint.Tags))
	for i, tag := range endpoint.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &EndpointList{
		Id:           endpoint.Id,
		ProgramId:    endpoint.ProgramId,
		ProgramName:  programName,
		Domain:       endpoint.Domain,
		URI:          endpoint.URI,
		Method:       endpoint.Method,
		EndpointType: string(endpoint.EndpointType),
		CreatedAt:    endpoint.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    endpoint.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Text:         text,
		Tags:         tags,
	}
}

type EndpointDetail struct {
	Id           int           `json:"id"`
	ProgramId    int           `json:"program_id"`
	ProgramName  string        `json:"program_name"`
	Domain       string        `json:"domain"`
	URI          string        `json:"uri"`
	Method       string        `json:"method"`
	EndpointType string        `json:"endpoint_type"`
	Description  string        `json:"description"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	Notes        []NoteListing `json:"notes"`
	Attachments  []Attachment  `json:"attachments"`
	Images       []Image       `json:"images"`
	Tags         []TagDTO      `json:"tags"`
}

func ToEndpointDetail(endpoint *models.Endpoint) *EndpointDetail {
	programName := ""
	if endpoint.Program != nil {
		programName = endpoint.Program.Name
	}

	notes := make([]NoteListing, len(endpoint.Notes))
	for i, note := range endpoint.Notes {
		notePtr := &note
		notes[i] = *ToNoteListing(notePtr)
	}

	attachments := make([]Attachment, len(endpoint.Attachments))
	for i, attachment := range endpoint.Attachments {
		attachments[i] = Attachment{
			Id:           attachment.Id,
			Filename:     attachment.Filename,
			OriginalName: attachment.OriginalName,
			URL:          attachment.GetURL(),
		}
	}

	images := make([]Image, len(endpoint.Images))
	for i, image := range endpoint.Images {
		images[i] = Image{
			Id:           image.Id,
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.GetURL(),
		}
	}

	tags := make([]TagDTO, len(endpoint.Tags))
	for i, tag := range endpoint.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &EndpointDetail{
		Id:           endpoint.Id,
		ProgramId:    endpoint.ProgramId,
		ProgramName:  programName,
		Domain:       endpoint.Domain,
		URI:          endpoint.URI,
		Method:       endpoint.Method,
		EndpointType: string(endpoint.EndpointType),
		Description:  endpoint.Note,
		CreatedAt:    endpoint.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    endpoint.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Notes:        notes,
		Attachments:  attachments,
		Images:       images,
		Tags:         tags,
	}
}

// ===== Requests =====
type RequestList struct {
	Id               int      `json:"id"`
	ProgramId        int      `json:"program_id"`
	ProgramName      string   `json:"program_name"`
	EndpointId       int      `json:"endpoint_id"`
	EndpointName     string   `json:"endpoint_name"`
	JobId            int      `json:"job_id"`
	SequenceNumber   int      `json:"sequence_number"`
	URL              string   `json:"url"`
	Method           string   `json:"method"`
	Domain           string   `json:"domain"`
	StatusCode       int      `json:"status_code"`
	ContentType      string   `json:"content_type"`
	Size             int      `json:"size"`
	ReqHash          string   `json:"req_hash"`
	ResponseHash     string   `json:"response_hash"`
	ResponseBodyHash string   `json:"response_body_hash"`
	Text             string   `json:"text"`
	Tags             []TagDTO `json:"tags"`
}

type RequestDetail struct {
	Id               int           `json:"id"`
	ProgramId        int           `json:"program_id"`
	ProgramName      string        `json:"program_name"`
	EndpointId       int           `json:"endpoint_id"`
	EndpointName     string        `json:"endpoint_name"`
	JobId            int           `json:"job_id"`
	SequenceNumber   int           `json:"sequence_number"`
	URL              string        `json:"url"`
	Method           string        `json:"method"`
	Domain           string        `json:"domain"`
	StatusCode       int           `json:"status_code"`
	RequestHeaders   string        `json:"request_headers"`
	RequestBody      interface{}   `json:"request_body"`
	ResponseBody     interface{}   `json:"response_body"`
	ResponseHeaders  interface{}   `json:"response_headers"`
	ReqHash          string        `json:"req_hash"`
	ResponseHash     string        `json:"response_hash"`
	ResponseBodyHash string        `json:"response_body_hash"`
	LatencyMs        int           `json:"latency_ms"`
	Notes            []NoteListing `json:"notes"`
	Attachments      []Attachment  `json:"attachments"`
	Images           []Image       `json:"images"`
	Tags             []TagDTO      `json:"tags"`
}

// ===== Jobs =====
type Job struct {
	Id          int    `json:"id"`
	JobType     string `json:"job_type" validate:"oneof=import_har import_burp_xml"`
	Title       string `json:"title"`
	Progress    int    `json:"progress" validate:"min=1,max=100"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}

// ===== Import HAR =====
type ImportHarRequest struct {
	ProgramId      int    `form:"program_id" validate:"required"`
	IgnoredHeaders string `form:"ignored_headers"`
	// file handled as multipart
}

// ===== Request Conversion Functions =====

// extractContentType extracts content type from response headers
func extractContentType(responseHeaders string) string {
	if responseHeaders == "" {
		return ""
	}

	// Parse the JSON headers
	var headers map[string]interface{}
	if err := json.Unmarshal([]byte(responseHeaders), &headers); err != nil {
		// If JSON parsing fails, try to parse as HeaderSlice
		if headerSlice, err := models.HeaderSliceFromJSON(responseHeaders); err == nil {
			for _, header := range headerSlice {
				if strings.EqualFold(header.Name, "content-type") {
					return header.Value
				}
			}
		}
		return ""
	}

	// Look for Content-Type header (case insensitive)
	for key, value := range headers {
		if strings.EqualFold(key, "content-type") {
			if str, ok := value.(string); ok {
				return str
			}
		}
	}

	return ""
}

func ToRequestList(request *models.MyRequest) *RequestList {
	programName := ""
	programId := 0
	if request.ProgramId != nil {
		programId = *request.ProgramId
		if request.Program != nil {
			programName = request.Program.Name
		}
	}

	endpointName := ""
	if request.Endpoint != nil {
		endpointName = request.Endpoint.URI
	}

	// Parse headers and body for text concatenation
	_, _ = services.ParseRequestHeaders(request.ReqHeaders)
	var responseHeaders interface{}
	if request.ResHeaders == "" {
		responseHeaders = map[string]interface{}{}
	} else {
		parsedHeaders, err := services.ParseResponseHeaders(request.ResHeaders)
		if err != nil {
			responseHeaders = request.ResHeaders
		} else {
			responseHeaders = parsedHeaders
		}
	}
	requestBody, _ := services.ParseRequestBody(request.ReqBody)
	responseBody, _ := services.ParseResponseBody(request.ResBody)

	// Concatenate all information into text field
	text := fmt.Sprintf("ID: %d\nProgram ID: %d\nProgram Name: %s\nEndpoint ID: %d\nEndpoint Name: %s\nJob ID: %d\nSequence Number: %d\nURL: %s\nMethod: %s\nDomain: %s\nStatus Code: %d\nContent Type: %s\nSize: %d\nRequest Headers: %s\nRequest Body: %v\nResponse Headers: %v\nResponse Body: %v\nRequest Hash: %s\nResponse Hash: %s\nResponse Body Hash: %s",
		request.Id,
		programId,
		programName,
		request.EndpointId,
		endpointName,
		request.ImportJobId,
		request.Sequence,
		request.URL,
		request.Method,
		request.Domain,
		request.ResStatus,
		extractContentType(request.ResHeaders),
		request.RespSize,
		request.ReqHeaders,
		requestBody,
		responseHeaders,
		responseBody,
		request.ReqHash,
		request.ResHash,
		request.ResBodyHash,
	)

	// Add notes to text
	if len(request.Notes) > 0 {
		text += "\nNotes:\n"
		for _, note := range request.Notes {
			text += fmt.Sprintf("- %s\n", note.Value)
		}
	}

	// Add attachments to text
	if len(request.Attachments) > 0 {
		text += "\nAttachments:\n"
		for _, attachment := range request.Attachments {
			text += fmt.Sprintf("- %s (%s)\n", attachment.OriginalName, attachment.Filename)
		}
	}

	// Convert tags
	tags := make([]TagDTO, len(request.Tags))
	for i, tag := range request.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &RequestList{
		Id:               request.Id,
		ProgramId:        programId,
		ProgramName:      programName,
		EndpointId:       request.EndpointId,
		EndpointName:     endpointName,
		JobId:            request.ImportJobId,
		SequenceNumber:   request.Sequence,
		URL:              request.URL,
		Method:           request.Method,
		Domain:           request.Domain,
		StatusCode:       request.ResStatus,
		ContentType:      extractContentType(request.ResHeaders),
		Size:             request.RespSize,
		ReqHash:          request.ReqHash,
		ResponseHash:     request.ResHash,
		ResponseBodyHash: request.ResBodyHash,
		Text:             text,
		Tags:             tags,
	}
}

func ToRequestDetail(request *models.MyRequest) *RequestDetail {
	programName := ""
	programId := 0
	if request.ProgramId != nil {
		programId = *request.ProgramId
		if request.Program != nil {
			programName = request.Program.Name
		}
	}

	endpointName := ""
	if request.Endpoint != nil {
		endpointName = request.Endpoint.URI
	}

	// Parse headers and body
	_, _ = services.ParseRequestHeaders(request.ReqHeaders)
	var responseHeaders interface{}
	if request.ResHeaders == "" {
		// If response headers are empty, return empty object
		responseHeaders = map[string]interface{}{}
	} else {
		parsedHeaders, err := services.ParseResponseHeaders(request.ResHeaders)
		if err != nil {
			// If parsing fails, return the raw headers as string
			responseHeaders = request.ResHeaders
		} else {
			responseHeaders = parsedHeaders
		}
	}
	requestBody, _ := services.ParseRequestBody(request.ReqBody)
	responseBody, _ := services.ParseResponseBody(request.ResBody)

	// Convert notes
	notes := make([]NoteListing, len(request.Notes))
	for i, note := range request.Notes {
		notePtr := &note
		notes[i] = *ToNoteListing(notePtr)
	}

	// Convert attachments
	attachments := make([]Attachment, len(request.Attachments))
	for i, attachment := range request.Attachments {
		attachments[i] = Attachment{
			Id:           attachment.Id,
			Filename:     attachment.Filename,
			OriginalName: attachment.OriginalName,
			URL:          attachment.GetURL(),
		}
	}

	// Convert images
	images := make([]Image, len(request.Images))
	for i, image := range request.Images {
		images[i] = Image{
			Id:           image.Id,
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.GetURL(),
		}
	}

	// Convert tags
	tags := make([]TagDTO, len(request.Tags))
	for i, tag := range request.Tags {
		tags[i] = *ToTagDTO(&tag)
	}

	return &RequestDetail{
		Id:               request.Id,
		ProgramId:        programId,
		ProgramName:      programName,
		EndpointId:       request.EndpointId,
		EndpointName:     endpointName,
		JobId:            request.ImportJobId,
		SequenceNumber:   request.Sequence,
		URL:              request.URL,
		Method:           request.Method,
		Domain:           request.Domain,
		StatusCode:       request.ResStatus,
		RequestHeaders:   request.ReqHeaders,
		RequestBody:      requestBody,
		ResponseBody:     responseBody,
		ResponseHeaders:  responseHeaders,
		ReqHash:          request.ReqHash,
		ResponseHash:     request.ResHash,
		ResponseBodyHash: request.ResBodyHash,
		LatencyMs:        int(request.LatencyMs),
		Notes:            notes,
		Attachments:      attachments,
		Images:           images,
		Tags:             tags,
	}
}

// ===== Tags =====
type TagInput struct {
	Name     string `json:"name" validate:"required"`
	Priority int    `json:"priority"`
}

type TagDTO struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func (input *TagInput) ToModel() *models.Tag {
	return &models.Tag{
		Name:     input.Name,
		Priority: input.Priority,
	}
}

func ToTagDTO(tag *models.Tag) *TagDTO {
	return &TagDTO{
		Id:       tag.Id,
		Name:     tag.Name,
		Priority: tag.Priority,
	}
}
