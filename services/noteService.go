package services

import (
	"context"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type NoteService struct {
	DB *gorm.DB
}

// NewInstance returns a copy of the service with a new transaction
func (s *NoteService) NewInstance(ctx context.Context) (*NoteService, func(), func() error) {
	tx := s.DB.WithContext(ctx).Begin()
	return &NoteService{
			DB: tx,
		}, func() {
			tx.Rollback()
		}, func() error {
			return tx.Commit().Error
		}
}

// CloneWithDb creates a new service instance with the given DB
func (s *NoteService) CloneWithDb(db *gorm.DB) *NoteService {
	return &NoteService{
		DB: db,
	}
}

func (s *NoteService) validate(db *gorm.DB, id int, input *models.Note) error {
	return nil
}

// Create creates a new note and returns its Id
func (s *NoteService) Create(ctx context.Context, note *models.Note) (int, error) {
	if err := s.validate(s.DB.WithContext(ctx), 0, note); err != nil {
		return 0, err
	}
	if err := s.DB.WithContext(ctx).Create(note).Error; err != nil {
		return 0, err
	}
	return note.Id, nil
}

// Get retrieves a note by Id
func (s *NoteService) Get(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note
	if err := s.DB.WithContext(ctx).Preload("Taggables.Tag").First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// List retrieves all notes with optional filtering
func (s *NoteService) List(ctx context.Context, referenceType string, search string) ([]*models.Note, error) {
	var notes []*models.Note
	query := s.DB.WithContext(ctx)

	if referenceType != "" {
		query = query.Where("reference_type = ?", referenceType)
	}

	if search != "" {
		query = query.Where("value LIKE ?", "%"+search+"%")
	}

	if err := query.Preload("Taggables.Tag").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// Update updates a note's value and returns its Id
func (s *NoteService) Update(ctx context.Context, id int, value string) (int, error) {
	note, err := first[models.Note](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}

	updates := map[string]any{
		"Value": value,
	}
	if err := s.DB.WithContext(ctx).Model(&note).Updates(updates).Error; err != nil {
		return 0, err
	}
	return note.Id, nil
}

// Delete deletes a note by Id and returns the deleted Id
func (s *NoteService) Delete(ctx context.Context, id int) (int, error) {
	note, err := first[models.Note](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Delete the note
	if err := tx.Delete(&note).Error; err != nil {
		return 0, err
	}

	// Clean up related dependencies
	err = tx.Exec("DELETE FROM taggables WHERE taggable_type = ? AND taggable_id = ?", models.TaggableTypeNotes, id).Error
	if err != nil {
		return 0, err
	}

	return note.Id, tx.Commit().Error
}
