package models

import (
	"time"
)

// Attachment represents a file attachment linked to a program, endpoint, or request
type Attachment struct {
	Id            int       `gorm:"primaryKey"`
	Filename      string    `gorm:"size:255;not null"`
	OriginalName  string    `gorm:"size:255;not null"`
	FilePath      string    `gorm:"size:500;not null"`
	FileSize      int64     `gorm:"not null"`
	MimeType      string    `gorm:"size:100"`
	ReferenceType string    `gorm:"size:20;not null;index"` // "programs", "endpoints", "requests"
	ReferenceID   int       `gorm:"not null;index"`         // Changed from ReferenceId to ReferenceID for GORM polymorphic
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// GetURL returns the URL for accessing the attachment
func (a *Attachment) GetURL() string {
	return "/attachments/" + a.Filename
}
