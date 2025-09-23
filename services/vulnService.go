package services

import (
	"context"
	"fmt"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type VulnService struct {
	DB *gorm.DB
}

// Create creates a new vulnerability and returns its Id
func (s *VulnService) Create(ctx context.Context, vuln *models.Vuln) (int, error) {
	// Generate slug if not provided
	if vuln.Slug == "" {
		vuln.Slug = vuln.GenerateSlug()
	}

	// Validate parent exists if ParentId is provided
	if vuln.ParentId != nil {
		var parentVuln models.Vuln
		if err := s.DB.WithContext(ctx).First(&parentVuln, *vuln.ParentId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return 0, fmt.Errorf("parent vulnerability with ID %d not found", *vuln.ParentId)
			}
			return 0, fmt.Errorf("failed to validate parent vulnerability: %v", err)
		}
	}

	if err := s.DB.WithContext(ctx).Create(vuln).Error; err != nil {
		return 0, fmt.Errorf("failed to create vulnerability: %v", err)
	}

	return vuln.Id, nil
}

// Get retrieves a vulnerability by ID with all associations
func (s *VulnService) Get(ctx context.Context, id int) (*models.Vuln, error) {
	var vuln models.Vuln
	if err := s.DB.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Preload("Attachments").
		Preload("Images").
		Preload("Notes").
		First(&vuln, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vulnerability not found")
		}
		return nil, fmt.Errorf("failed to get vulnerability: %v", err)
	}
	return &vuln, nil
}

// List retrieves all vulnerabilities with optional filtering
func (s *VulnService) List(ctx context.Context, parentId *int) ([]*models.Vuln, error) {
	var vulns []*models.Vuln
	query := s.DB.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Preload("Attachments").
		Preload("Images").
		Preload("Notes")

	if parentId != nil {
		query = query.Where("parent_id = ?", *parentId)
	}

	if err := query.Find(&vulns).Error; err != nil {
		return nil, fmt.Errorf("failed to list vulnerabilities: %v", err)
	}

	return vulns, nil
}

// Update updates an existing vulnerability and returns its Id
func (s *VulnService) Update(ctx context.Context, id int, vuln *models.Vuln) (int, error) {
	// Check if vulnerability exists
	var existingVuln models.Vuln
	if err := s.DB.WithContext(ctx).First(&existingVuln, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("vulnerability not found")
		}
		return 0, fmt.Errorf("failed to find vulnerability: %v", err)
	}

	// Validate parent exists if ParentId is provided and different from current
	if vuln.ParentId != nil && *vuln.ParentId != id {
		var parentVuln models.Vuln
		if err := s.DB.WithContext(ctx).First(&parentVuln, *vuln.ParentId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return 0, fmt.Errorf("parent vulnerability with ID %d not found", *vuln.ParentId)
			}
			return 0, fmt.Errorf("failed to validate parent vulnerability: %v", err)
		}
	}

	// Update fields
	existingVuln.Title = vuln.Title
	existingVuln.Body = vuln.Body
	existingVuln.ParentId = vuln.ParentId

	// Generate new slug if title changed
	if existingVuln.Title != vuln.Title {
		existingVuln.Slug = existingVuln.GenerateSlug()
	}

	if err := s.DB.WithContext(ctx).Save(&existingVuln).Error; err != nil {
		return 0, fmt.Errorf("failed to update vulnerability: %v", err)
	}

	return existingVuln.Id, nil
}

// Delete deletes a vulnerability by ID and returns the deleted Id
func (s *VulnService) Delete(ctx context.Context, id int) (int, error) {
	// Check if vulnerability exists
	var vuln models.Vuln
	if err := s.DB.WithContext(ctx).First(&vuln, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("vulnerability not found")
		}
		return 0, fmt.Errorf("failed to find vulnerability: %v", err)
	}

	// Check if vulnerability has children
	var childCount int64
	if err := s.DB.WithContext(ctx).Model(&models.Vuln{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return 0, fmt.Errorf("failed to check for child vulnerabilities: %v", err)
	}

	if childCount > 0 {
		return 0, fmt.Errorf("cannot delete vulnerability with child vulnerabilities")
	}

	// Delete the vulnerability
	if err := s.DB.WithContext(ctx).Delete(&vuln).Error; err != nil {
		return 0, fmt.Errorf("failed to delete vulnerability: %v", err)
	}

	return vuln.Id, nil
}

// GetBySlug retrieves a vulnerability by slug
func (s *VulnService) GetBySlug(ctx context.Context, slug string) (*models.Vuln, error) {
	var vuln models.Vuln
	if err := s.DB.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Preload("Attachments").
		Preload("Images").
		Preload("Notes").
		Where("slug = ?", slug).
		First(&vuln).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vulnerability not found")
		}
		return nil, fmt.Errorf("failed to get vulnerability by slug: %v", err)
	}
	return &vuln, nil
}
