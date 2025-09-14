package models

import "time"

type ImportJob struct {
	Id             int       `gorm:"primaryKey"`
	ProgramId      *int      `gorm:"index"`            // Foreign key to Program (nullable for migration)
	JobType        string    `gorm:"size:20;not null"` // "import_har", "import_xml"
	Title          string    `gorm:"not null"`
	Progress       int       `gorm:"not null;default:0"` // 0-100
	Description    string    `gorm:"type:text"`
	IgnoredHeaders string    `gorm:"type:text"` // Store as JSON string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	// One-to-many relationship
	Requests []MyRequest `gorm:"foreignKey:ImportJobId"`
}
