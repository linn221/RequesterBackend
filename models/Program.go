package models

import "time"

// Program represents a program/project
type Program struct {
	Id        int       `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	URL       string    `gorm:"size:500"`
	Note      string    `gorm:"type:text"`
	Scope     string    `gorm:"type:text"`
	Domains   string    `gorm:"type:text"` // Store as JSON string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// One-to-many relationships
	// ImportJobs []ImportJob `gorm:"foreignKey:ProgramId"`
	Endpoints []Endpoint  `gorm:"foreignKey:ProgramId"`
	Requests  []MyRequest `gorm:"foreignKey:ProgramId"`

	// Polymorphic relationships
	Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	Notes       []Note       `gorm:"polymorphic:Reference;polymorphicValue:programs"`
}
