package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type MyRequest struct {
	Id          int    `gorm:"primaryKey"`
	ProgramId   *int   `gorm:"index"`          // Foreign key to Program (nullable for migration)
	ImportJobId int    `gorm:"not null;index"` // Foreign key to ImportJob
	EndpointId  int    `gorm:"not null;index"` // Foreign key to Endpoint
	Sequence    int    `gorm:"not null"`
	URL         string `gorm:"type:text;not null"`
	Method      string `gorm:"size:10;not null"`
	Domain      string `gorm:"size:255;not null"`

	ReqHeaders string `gorm:"type:text"` // Store as JSON string
	ReqBody    string `gorm:"type:longtext"`

	ResStatus  int    `gorm:"not null"`
	ResHeaders string `gorm:"type:text"` // Store as JSON string
	ResBody    string `gorm:"type:longtext"`
	RespSize   int    `gorm:"not null"`
	LatencyMs  int64  `gorm:"not null"`

	RequestTime string `gorm:"size:50"`
	// hashes
	ReqHash1    string `gorm:"size:64;index"` // hash raw request
	ReqHash     string `gorm:"size:64;index"`
	ResHash     string `gorm:"size:64;index"`
	ResBodyHash string `gorm:"size:64;index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Belongs to relationships
	Program  *Program  `gorm:"foreignKey:ProgramId"`
	Endpoint *Endpoint `gorm:"foreignKey:EndpointId"`

	// Polymorphic relationships
	Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:requests"`
	Notes       []Note       `gorm:"polymorphic:Reference;polymorphicValue:requests"`
	Images      []Image      `gorm:"polymorphic:Reference;polymorphicValue:requests"`
}

// Temporary struct for parsing HAR files (with HeaderSlice fields)
type TempMyRequest struct {
	Sequence    int
	URL         string
	Method      string
	Domain      string
	ReqHeaders  HeaderSlice
	ReqBody     string
	ResStatus   int
	ResHeaders  HeaderSlice
	ResBody     string
	RespSize    int
	LatencyMs   int64
	RequestTime string
	// hashes
	ReqHash1    string
	ReqHash     string
	ResHash     string
	ResBodyHash string
}

type HeaderSlice []Header

func (hs HeaderSlice) EchoAll() string {
	var result string
	for _, h := range hs {
		result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
	}
	return result
}

func (hs HeaderSlice) EchoMatcher(headerNames ...string) string {
	matchMap := make(map[string]struct{}, len(headerNames))
	for _, hname := range headerNames {
		matchMap[strings.ToLower(hname)] = struct{}{}
	}

	var result string
	for _, h := range hs {
		if _, match := matchMap[strings.ToLower(h.Name)]; match {
			result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
		}
	}
	return result
}

func (hs HeaderSlice) EchoFilter(headerNames ...string) string {
	filterMap := make(map[string]struct{}, len(headerNames))
	for _, hname := range headerNames {
		filterMap[strings.ToLower(hname)] = struct{}{}
	}

	var result string
	for _, h := range hs {
		if _, filter := filterMap[strings.ToLower(h.Name)]; !filter {
			result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
		}
	}
	return result
}

// Convert HeaderSlice to JSON string for database storage
func (hs HeaderSlice) ToJSON() (string, error) {
	data, err := json.Marshal(hs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Convert JSON string to HeaderSlice from database
func HeaderSliceFromJSON(jsonStr string) (HeaderSlice, error) {
	var hs HeaderSlice
	err := json.Unmarshal([]byte(jsonStr), &hs)
	return hs, err
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s\n", h.Name, h.Value)
}

func (r TempMyRequest) requestText() string {
	raw := r.Method + " " + r.URL + " " + r.ReqBody + " " + r.ReqHeaders.EchoAll()
	return raw
}

// Convert TempMyRequest to MyRequest for database storage
func (temp *TempMyRequest) ToMyRequest(programId, importJobId, endpointId int) (*MyRequest, error) {
	// Convert HeaderSlice to JSON strings
	reqHeadersJSON, err := temp.ReqHeaders.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert request headers to JSON: %v", err)
	}

	resHeadersJSON, err := temp.ResHeaders.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert response headers to JSON: %v", err)
	}

	return &MyRequest{
		ProgramId:   &programId,
		ImportJobId: importJobId,
		EndpointId:  endpointId,
		Sequence:    temp.Sequence,
		URL:         temp.URL,
		Method:      temp.Method,
		Domain:      temp.Domain,
		ReqHeaders:  reqHeadersJSON,
		ReqBody:     temp.ReqBody,
		ResStatus:   temp.ResStatus,
		ResHeaders:  resHeadersJSON,
		ResBody:     temp.ResBody,
		RespSize:    temp.RespSize,
		LatencyMs:   temp.LatencyMs,
		RequestTime: temp.RequestTime,
		ReqHash1:    temp.ReqHash1,
		ReqHash:     temp.ReqHash,
		ResHash:     temp.ResHash,
		ResBodyHash: temp.ResBodyHash,
	}, nil
}
