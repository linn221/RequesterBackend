package models

import "time"

// EndpointType represents the type of endpoint
type EndpointType string

const (
	EndpointTypeWeb     EndpointType = "Web"
	EndpointTypeAPI     EndpointType = "API"
	EndpointTypeGraphQL EndpointType = "GraphQL"
)

// Endpoint represents an API endpoint
type Endpoint struct {
	Id           int          `gorm:"primaryKey"`
	ProgramId    int          `gorm:"index;not null"` // Foreign key to Program
	Method       string       `gorm:"size:10;not null"`
	Domain       string       `gorm:"size:255;not null"`
	URI          string       `gorm:"type:text;not null"`
	EndpointType EndpointType `gorm:"size:20;not null;default:'API'"`
	Note         string       `gorm:"type:text"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`

	// Belongs to relationship
	Program *Program `gorm:"foreignKey:ProgramId"`

	// One-to-many relationship
	Requests []MyRequest `gorm:"foreignKey:EndpointId"`

	// Polymorphic relationships
	Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:endpoints"`
	Notes       []Note       `gorm:"polymorphic:Reference;polymorphicValue:endpoints"`
}
