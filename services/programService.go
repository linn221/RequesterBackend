package services

import (
	"context"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type ProgramService struct {
	DB *gorm.DB
}

// NewInstance returns a copy of the service with a new transaction
func (s *ProgramService) NewInstance(ctx context.Context) (*ProgramService, func(), func() error) {
	tx := s.DB.WithContext(ctx).Begin()
	return &ProgramService{
			DB: tx,
		}, func() {
			tx.Rollback()
		}, func() error {
			return tx.Commit().Error
		}
}

// CloneWithDb creates a new service instance with the given DB
func (s *ProgramService) CloneWithDb(db *gorm.DB) *ProgramService {
	return &ProgramService{
		DB: db,
	}
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
	if err := s.DB.WithContext(ctx).Preload("Endpoints").Preload("Requests").Preload("Notes").Preload("Taggables.Tag").First(&program, id).Error; err != nil {
		return nil, err
	}
	return &program, nil
}

// List retrieves all programs
func (s *ProgramService) List(ctx context.Context) ([]*models.Program, error) {
	var programs []*models.Program
	if err := s.DB.WithContext(ctx).Preload("Endpoints").Preload("Requests").Preload("Notes").Preload("Taggables.Tag").Find(&programs).Error; err != nil {
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

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Delete the program
	if err := tx.Delete(&program).Error; err != nil {
		return 0, err
	}

	// Clean up related dependencies
	err = tx.Exec("DELETE FROM taggables WHERE taggable_type = ? AND taggable_id = ?", models.TaggableTypePrograms, id).Error
	if err != nil {
		return 0, err
	}
	err = tx.Exec("DELETE FROM notes WHERE reference_type = 'programs' AND reference_id = ?", id).Error
	if err != nil {
		return 0, err
	}

	return program.Id, tx.Commit().Error
}
