package services

import (
	"context"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type EndpointService struct {
	DB *gorm.DB
}

func (s *EndpointService) validate(db *gorm.DB, id int, input *models.Endpoint) error {
	return nil
}

// Create creates a new endpoint and returns its Id
func (s *EndpointService) Create(ctx context.Context, endpoint *models.Endpoint) (int, error) {
	if err := s.validate(s.DB.WithContext(ctx), 0, endpoint); err != nil {
		return 0, err
	}
	if err := s.DB.WithContext(ctx).Create(endpoint).Error; err != nil {
		return 0, err
	}
	return endpoint.Id, nil
}

// Get retrieves an endpoint by Id
func (s *EndpointService) Get(ctx context.Context, id int) (*models.Endpoint, error) {
	var endpoint models.Endpoint
	if err := s.DB.WithContext(ctx).Preload("Program").Preload("Notes").First(&endpoint, id).Error; err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// List retrieves all endpoints
func (s *EndpointService) List(ctx context.Context) ([]*models.Endpoint, error) {
	var endpoints []*models.Endpoint
	if err := s.DB.WithContext(ctx).Preload("Program").Preload("Notes").Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

// Update updates an existing endpoint and returns its Id
func (s *EndpointService) Update(ctx context.Context, id int, input *models.Endpoint) (int, error) {

	if err := s.validate(s.DB.WithContext(ctx), id, input); err != nil {
		return 0, err
	}

	endpoint, err := first[models.Endpoint](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}
	updates := map[string]any{ // from EndpointInput
		"ProgramId":    input.ProgramId,
		"Method":       input.Method,
		"Domain":       input.Domain,
		"URI":          input.URI,
		"EndpointType": input.EndpointType,
		"Note":         input.Note,
	}
	if err := s.DB.WithContext(ctx).Model(&endpoint).Updates(updates).Error; err != nil {
		return 0, err
	}
	return endpoint.Id, nil
}

// Delete deletes an endpoint by Id and returns the deleted Id
func (s *EndpointService) Delete(ctx context.Context, id int) (int, error) {
	endpoint, err := first[models.Endpoint](s.DB.WithContext(ctx), id)
	if err != nil {
		return 0, err
	}

	if err := s.DB.WithContext(ctx).Delete(&endpoint).Error; err != nil {
		return 0, err
	}
	return endpoint.Id, nil
}
