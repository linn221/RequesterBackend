package models

import (
	"time"
)

// TaggableType represents the type of taggable resource
type TaggableType string

const (
	TaggableTypePrograms  TaggableType = "programs"
	TaggableTypeEndpoints TaggableType = "endpoints"
	TaggableTypeRequests  TaggableType = "requests"
	TaggableTypeVulns     TaggableType = "vulns"
	TaggableTypeNotes     TaggableType = "notes"
)

// Tag represents a tag that can be applied to various resources
type Tag struct {
	Id        int       `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null;uniqueIndex"`
	Priority  int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Many-to-many polymorphic relationships through Taggable
	Taggables []Taggable `gorm:"foreignKey:TagID"`
}

// Taggable represents the many-to-many polymorphic relationship between tags and resources
type Taggable struct {
	ID           int       `gorm:"primaryKey"`
	TagID        int       `gorm:"column:tag_id;not null;index"`
	TaggableType string    `gorm:"column:taggable_type;size:20;not null;index"` // "programs", "endpoints", "requests", "vulns", "notes"
	TaggableID   int       `gorm:"column:taggable_id;not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relationships
	Tag Tag `gorm:"foreignKey:TagID"`
}
