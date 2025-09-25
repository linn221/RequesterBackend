package models

import "time"

// Note represents a note linked to a program, endpoint, or request
type Note struct {
	Id            int       `gorm:"primaryKey"`
	ReferenceType string    `gorm:"size:20;not null;index"` // "programs", "endpoints", "requests"
	ReferenceID   int       `gorm:"not null;index"`         // Changed from ReferenceId to ReferenceID for GORM polymorphic
	Value         string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	// Polymorphic relationships
	Taggables []Taggable `gorm:"polymorphic:Taggable;polymorphicValue:notes"`
}
