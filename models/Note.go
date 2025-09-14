package models

import "time"

// Note represents a note linked to a program, endpoint, or request
type Note struct {
	Id            int       `gorm:"primaryKey"`
	ReferenceType string    `gorm:"size:20;not null;index"` // "programs", "endpoints", "requests"
	ReferenceId   int       `gorm:"not null;index"`
	Value         string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
