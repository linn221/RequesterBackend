package handlers

import "github.com/linn221/RequesterBackend/models"

// ===== Notes =====
type CreateNoteRequest struct {
	ReferenceType string `json:"reference_type" validate:"required,oneof=programs endpoints requests"`
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
	Id            int    `json:"id"`
	ReferenceType string `json:"reference_type"`
	ReferenceId   int    `json:"reference_id"`
	Value         string `json:"value"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func ToNote(note *models.Note) *Note {
	return &Note{
		Id:            note.Id,
		ReferenceType: note.ReferenceType,
		ReferenceId:   note.ReferenceId,
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
	return &NoteDetail{
		Id:            note.Id,
		ReferenceType: note.ReferenceType,
		ReferenceId:   note.ReferenceId,
		Value:         note.Value,
		CreatedAt:     note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ===== Attachments =====
type Attachment struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

type UploadAttachmentRequest struct {
	ReferenceType string `json:"reference_type" validate:"required,oneof=programs endpoints requests"`
	ReferenceId   int    `json:"reference_id" validate:"required"`
	// file is handled separately as multipart
}

// ===== Programs =====
type ProgramInput struct {
	Name    string `json:"name" validate:"required"`
	URL     string `json:"url" validate:"required,url"`
	Scope   string `json:"scope" validate:"required"`
	Domains string `json:"domains" validate:"required"`
	Note    string `json:"note" validate:"required"`
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
	Id   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func ToProgramList(program *models.Program) *ProgramList {
	return &ProgramList{
		Id:   program.Id,
		Name: program.Name,
		URL:  program.URL,
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
			Id:       attachment.Id,
			Filename: attachment.Filename,
			URL:      attachment.GetURL(),
		}
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
	Id           int    `json:"id"`
	ProgramId    int    `json:"program_id"`
	ProgramName  string `json:"program_name"`
	Domain       string `json:"domain"`
	URI          string `json:"uri"`
	Method       string `json:"method"`
	EndpointType string `json:"endpoint_type"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func ToEndpointList(endpoint *models.Endpoint) *EndpointList {
	programName := ""
	if endpoint.Program != nil {
		programName = endpoint.Program.Name
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
			Id:       attachment.Id,
			Filename: attachment.Filename,
			URL:      attachment.GetURL(),
		}
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
	}
}

// ===== Requests =====
type RequestList struct {
	Id             int      `json:"id"`
	ProgramId      int      `json:"program_id"`
	ProgramName    string   `json:"program_name"`
	EndpointId     int      `json:"endpoint_id"`
	EndpointName   string   `json:"endpoint_name"`
	JobId          int      `json:"job_id"`
	SequenceNumber int      `json:"sequence_number"`
	URL            string   `json:"url"`
	Method         string   `json:"method"`
	Domain         string   `json:"domain"`
	StatusCode     int      `json:"status_code"`
	SearchResults  []string `json:"search_results"`
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
	SearchResults    []string      `json:"search_results"`
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
}

// ===== Jobs =====
type Job struct {
	Id          int    `json:"id"`
	JobType     string `json:"job_type" validate:"oneof=import_har import_xml"`
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
