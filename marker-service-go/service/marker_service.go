package service

import (
	"context"
	"fmt"
	"time"

	"github.com/marker-service-go/models"
	"github.com/marker-service-go/repository"
)

// IMarkerService defines the interface for marker business logic
type IMarkerService interface {
	CreateMarkerAsync(ctx context.Context, userId int64, request *models.CreateMarkerRequest) (*models.MarkerResponse, error)
	GetMarkerAsync(ctx context.Context, userId int64, markerId string) (*models.MarkerResponse, error)
	GetAllMarkersAsync(ctx context.Context, userId int64, category *models.MarkerCategory) ([]*models.MarkerResponse, error)
	UpdateMarkerAsync(ctx context.Context, userId int64, markerId string, request *models.UpdateMarkerRequest) (*models.MarkerResponse, error)
	DeleteMarkerAsync(ctx context.Context, userId int64, markerId string) (bool, error)
}

// MarkerService implements IMarkerService
type MarkerService struct {
	repo repository.IMarkerRepository
}

// NewMarkerService creates a new marker service
func NewMarkerService(repo repository.IMarkerRepository) IMarkerService {
	return &MarkerService{
		repo: repo,
	}
}

// CreateMarkerAsync creates a new marker
func (s *MarkerService) CreateMarkerAsync(ctx context.Context, userId int64, request *models.CreateMarkerRequest) (*models.MarkerResponse, error) {
	// Validate input
	if request.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if !models.IsValidCategory(request.Category) {
		return nil, fmt.Errorf("invalid category")
	}

	if request.Latitude < -90 || request.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}

	if request.Longitude < -180 || request.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	// Generate marker ID
	now := time.Now().UTC()
	markerId := s.GenerateMarkerID(request.Category, now)

	// Create marker
	marker := &models.Marker{
		ID:          markerId,
		UserID:      userId,
		Name:        request.Name,
		Category:    request.Category,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Description: request.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Store in repository
	createdMarker, err := s.repo.CreateAsync(marker)
	if err != nil {
		return nil, err
	}

	return createdMarker.ToMarkerResponse(), nil
}

// GetMarkerAsync retrieves a specific marker
func (s *MarkerService) GetMarkerAsync(ctx context.Context, userId int64, markerId string) (*models.MarkerResponse, error) {
	marker, err := s.repo.GetByIdAsync(markerId, userId)
	if err != nil {
		return nil, err
	}

	if marker == nil {
		return nil, nil
	}

	return marker.ToMarkerResponse(), nil
}

// GetAllMarkersAsync retrieves all markers for a user, optionally filtered by category
func (s *MarkerService) GetAllMarkersAsync(ctx context.Context, userId int64, category *models.MarkerCategory) ([]*models.MarkerResponse, error) {
	markers, err := s.repo.GetByCategoryAsync(userId, category)
	if err != nil {
		return nil, err
	}

	var responses []*models.MarkerResponse
	for _, marker := range markers {
		responses = append(responses, marker.ToMarkerResponse())
	}

	return responses, nil
}

// UpdateMarkerAsync updates an existing marker
func (s *MarkerService) UpdateMarkerAsync(ctx context.Context, userId int64, markerId string, request *models.UpdateMarkerRequest) (*models.MarkerResponse, error) {
	// Get existing marker
	existing, err := s.repo.GetByIdAsync(markerId, userId)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, fmt.Errorf("marker not found")
	}

	// Apply updates only to provided fields
	if request.Name != nil {
		existing.Name = *request.Name
	}

	if request.Latitude != nil {
		if *request.Latitude < -90 || *request.Latitude > 90 {
			return nil, fmt.Errorf("latitude must be between -90 and 90")
		}
		existing.Latitude = *request.Latitude
	}

	if request.Longitude != nil {
		if *request.Longitude < -180 || *request.Longitude > 180 {
			return nil, fmt.Errorf("longitude must be between -180 and 180")
		}
		existing.Longitude = *request.Longitude
	}

	if request.Description != nil {
		existing.Description = request.Description
	}

	// Update timestamp
	existing.UpdatedAt = time.Now().UTC()

	// Save to repository
	updated, err := s.repo.UpdateAsync(existing)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, fmt.Errorf("marker not found")
	}

	return updated.ToMarkerResponse(), nil
}

// DeleteMarkerAsync removes a marker
func (s *MarkerService) DeleteMarkerAsync(ctx context.Context, userId int64, markerId string) (bool, error) {
	return s.repo.DeleteAsync(markerId, userId)
}

// GenerateMarkerID generates a marker ID in format {Category}_v1.0_{UnixTimestamp}
func (s *MarkerService) GenerateMarkerID(category models.MarkerCategory, timestamp time.Time) string {
	return fmt.Sprintf("%s_v1.0_%d", category.String(), timestamp.Unix())
}
