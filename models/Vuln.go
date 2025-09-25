package models

import (
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Vuln represents a vulnerability record
type Vuln struct {
	Id        int       `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Body      string    `gorm:"type:text;not null"`
	Slug      string    `gorm:"size:255;not null;uniqueIndex"`
	ParentId  *int      `gorm:"index"` // Self-referencing foreign key (nullable)
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Self-referencing relationship
	Parent   *Vuln   `gorm:"foreignKey:ParentId"`
	Children []*Vuln `gorm:"foreignKey:ParentId"`

	// Polymorphic relationships
	Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:vulns"`
	Images      []Image      `gorm:"polymorphic:Reference;polymorphicValue:vulns"`
	Notes       []Note       `gorm:"polymorphic:Reference;polymorphicValue:vulns"`
	Taggables   []Taggable   `gorm:"polymorphic:Taggable;polymorphicValue:vulns"`
}

// GenerateSlug creates a URL-friendly slug from the title
func (v *Vuln) GenerateSlug() string {
	// Convert to lowercase
	slug := strings.ToLower(v.Title)

	// Remove special characters and replace spaces with dashes
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace spaces and multiple dashes with single dash
	reg = regexp.MustCompile(`[\s-]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing dashes
	slug = strings.Trim(slug, "-")

	return slug
}

// BeforeCreate hook to generate slug before creating
func (v *Vuln) BeforeCreate(tx *gorm.DB) error {
	if v.Slug == "" {
		v.Slug = v.GenerateSlug()
	}
	return nil
}

// BeforeUpdate hook to regenerate slug if title changed
func (v *Vuln) BeforeUpdate(tx *gorm.DB) error {
	// Check if title has changed
	var oldVuln Vuln
	if err := tx.First(&oldVuln, v.Id).Error; err == nil {
		if oldVuln.Title != v.Title {
			v.Slug = v.GenerateSlug()
		}
	}
	return nil
}
