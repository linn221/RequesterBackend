package services

import (
	"context"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type ProgramService struct {
	DB *gorm.DB
}

func (s *ProgramService) validate(db *gorm.DB, id int, input *models.Program) error {
	return nil
}

// Create creates a new program and returns its Id
func (s *ProgramService) Create(ctx context.Context, program *models.Program) (int, error) {
	if err := s.validate(s.DB.WithContext(ctx), 0, program); err != nil {
		return 0, err
	}
	if err := s.DB.WithContext(ctx).Create(program).Error; err != nil {
		return 0, err
	}
	return program.Id, nil
}

// Get retrieves a program by Id
func (s *ProgramService) Get(ctx context.Context, id int) (*models.Program, error) {
	var program models.Program
	if err := s.DB.WithContext(ctx).Preload("Endpoints").Preload("Requests").Preload("Notes").First(&program, id).Error; err != nil {
		return nil, err
	}
	return &program, nil
}

// List retrieves all programs
func (s *ProgramService) List(ctx context.Context) ([]*models.Program, error) {
	var programs []*models.Program
	if err := s.DB.WithContext(ctx).Preload("Endpoints").Preload("Requests").Preload("Notes").Find(&programs).Error; err != nil {
		return nil, err
	}
	return programs, nil
}

// Update updates an existing program and returns its Id
func (s *ProgramService) Update(ctx context.Context, id int, input *models.Program) (int, error) {

	if err := s.validate(s.DB.WithContext(ctx), id, input); err != nil {
		return 0, err
	}

	program, err := first[models.Program](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}
	updates := map[string]any{ // from ProgramInput
		"Name":    input.Name,
		"URL":     input.URL,
		"Scope":   input.Scope,
		"Domains": input.Domains,
		"Note":    input.Note,
	}
	if err := s.DB.WithContext(ctx).Model(&program).Updates(updates).Error; err != nil {
		return 0, err
	}
	return program.Id, nil
}

// Delete deletes a program by Id and returns the deleted Id
func (s *ProgramService) Delete(ctx context.Context, id int) (int, error) {
	program, err := first[models.Program](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}

	if err := s.DB.WithContext(ctx).Delete(&program).Error; err != nil {
		return 0, err
	}
	return program.Id, nil
}
